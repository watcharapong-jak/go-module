package config

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"github.com/watcharapong-jak/go-module/fiberpkg/common"
	"gopkg.in/go-playground/validator.v9"
	"reflect"
	"regexp"
	"strings"
	"time"
	"unicode"
)

const (
	validatorRequired         = "required"
	validatorMin              = "min"
	validatorMax              = "max"
	validatorNotEqual         = "ne"
	validatorLen              = "len"
	validatorOneOf            = "oneof"
	validatorGreaterThan      = "gt"
	validatorGreaterThanEqual = "gte"
	validatorLessThan         = "lt"
	validatorLessThanEqual    = "lte"

	validateEmail       = "email"
	validateNumeric     = "numeric"
	validateDate        = "date"
	validateThai        = "thaialpha"
	validateAlphaNum    = "alphanum"
	validateNationality = "thaination"
	validateMoney       = "money"
	validateCurrency    = "currency"
	validateMobile      = "mobile"
)

var (
	mobileRegex  = regexp.MustCompile("^([0]{1}|(66))[0-9]{9,9}$")
	englishRegex = regexp.MustCompile("^[a-z,A-Z]*$")
	Validate     *validator.Validate
)

type LocaleDescription struct {
	Locale string
	TH     string
	EN     string
	CH     string
}

type ErrorCode struct {
	Code    string            `json:"code"`
	Message LocaleDescription `json:"description"`
}

func (ld LocaleDescription) MarshalJSON() ([]byte, error) {
	switch strings.ToLower(ld.Locale) {
	case "th":
		return json.Marshal(ld.TH)
	case "ch":
		return json.Marshal(ld.CH)
	default:
		return json.Marshal(ld.EN)
	}
}

var EM ErrorMessage

type ErrorMessage struct {
	vn         *viper.Viper
	ConfigPath string

	Success  ErrorCode
	Internal struct {
		General                 ErrorCode
		BadRequest              ErrorCode
		InternalServerError     ErrorCode
		DatabaseError           ErrorCode
		Timeout                 ErrorCode
		ParkingSlotNotAvailable ErrorCode
		LicensePlateInvalid     ErrorCode
		PickUpNumberNotUse      ErrorCode
		BookingDuplicate        ErrorCode
	}
	Validation struct {
		ValidationError           ErrorCode
		InvalidHashKey            ErrorCode
		DuplicateUsername         ErrorCode
		DuplicateTransactionRef   ErrorCode
		InsufficientAgentBalance  ErrorCode
		InsufficientPlayerBalance ErrorCode
		DecryptionFailure         ErrorCode
		InvalidRequestData        ErrorCode
		InvalidResponseData       ErrorCode
		HttpStatusNotOk           ErrorCode
		AgentSuspended            ErrorCode
		MemberSuspended           ErrorCode
		AlreadyCompletedCode      ErrorCode
		ApiReachLimitation        ErrorCode
		InactiveCurrency          ErrorCode
	}
	Feature struct {
		Auth struct {
			OTPInvalid       ErrorCode
			PhoneInvalid     ErrorCode
			PhoneDuplicated  ErrorCode
			SendSMS          ErrorCode
			MemberNotFound   ErrorCode
			PleaseVerifyOtp  ErrorCode
			AlreadyVerifyOtp ErrorCode
		}
	}
}

type ValidateStructResponse struct {
	Field string `json:"field"`
	Tag   string `json:"tag"`
	Value string `json:"value"`
}

func (em *ErrorMessage) Init() error {
	vn := viper.New()
	vn.AddConfigPath(em.ConfigPath)
	vn.SetConfigName("error")

	if err := vn.ReadInConfig(); err != nil {
		return err
	}

	em.vn = vn

	em.mapping("", reflect.ValueOf(em).Elem())

	return nil
}

func (ec ErrorCode) WithLocale(c *fiber.Ctx) ErrorCode {
	locale := c.AcceptsLanguages()
	if locale == "" {
		ec.Message.Locale = common.LangEnglish
	}
	ec.Message.Locale = locale
	return ec
}

func (em ErrorMessage) mapping(name string, v reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		fi := v.Field(i)
		if fi.Kind() != reflect.Struct {
			continue
		}

		fn := Underscore(v.Type().Field(i).Name)
		if name != "" {
			fn = fmt.Sprint(name, ".", fn)
		}

		if fi.Type().Name() == "ErrorCode" {
			fi.Set(reflect.ValueOf(em.ErrorCode(fn)))
			continue
		}
		em.mapping(fn, fi)
	}
}

func Underscore(str string) string {
	runes := []rune(str)
	var out []rune
	for i := 0; i < len(runes); i++ {
		if i > 0 && (unicode.IsUpper(runes[i]) || unicode.IsNumber(runes[i])) && ((i+1 < len(runes) && unicode.IsLower(runes[i+1])) || unicode.IsLower(runes[i-1])) {
			out = append(out, '_')
		}
		out = append(out, unicode.ToLower(runes[i]))
	}
	return string(out)
}

func (em ErrorMessage) ErrorCode(name string) ErrorCode {
	rtn := ErrorCode{
		Code: em.vn.GetString(fmt.Sprintf("%s.code", name)),
		Message: LocaleDescription{
			TH: em.vn.GetString(fmt.Sprintf("%s.th", name)),
			EN: em.vn.GetString(fmt.Sprintf("%s.en", name)),
			CH: em.vn.GetString(fmt.Sprintf("%s.ch", name)),
		},
	}
	return rtn
}

func ValidateStruct(v interface{}) ([]*ValidateStructResponse, error) {
	var errors []*ValidateStructResponse
	err := Validate.Struct(v)
	if err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var element ValidateStructResponse
			element.Field = err.Field()
			element.Tag = err.Tag()
			element.Value = err.Param()
			errors = append(errors, &element)
		}
	}
	return errors, err
}

func InitDefaultValidators() error {
	if Validate == nil {
		Validate = validator.New()
		if err := Validate.RegisterValidation("date", dateValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("thaialpha", thaiAlphaValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("englishalpha", thaiAlphaValidator); err != nil {
			return err
		}
		if err := Validate.RegisterValidation("mobile", mobileValidator); err != nil {
			return err
		}
		Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
	return nil
}

func dateValidator(fl validator.FieldLevel) bool {
	if _, err := time.Parse("2006/01/02", fl.Field().String()); err != nil {
		return false
	}
	return true
}

func thaiAlphaValidator(fl validator.FieldLevel) bool {
	s := []rune(fl.Field().String())

	for _, r := range s {
		if !unicode.Is(unicode.Thai, r) {
			return false
		}
	}
	return true
}

func engAlphaValidator(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return englishRegex.MatchString(v)
}

func mobileValidator(fl validator.FieldLevel) bool {
	v := fl.Field().String()
	if v == "" {
		return true
	}
	return mobileRegex.MatchString(v)
}
