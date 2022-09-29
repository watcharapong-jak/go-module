package appresponse

import (
	"github.com/goccy/go-json"
	"github.com/watcharapong-jak/go-module/fiberpkg/config"
)

type IResponse struct {
	config.ErrorCode
	Error           error       `json:"-"`
	ValidationError interface{} `json:"validationError,omitempty"`
	Data            interface{} `json:"data,omitempty"`
}

func (i IResponse) MarshalJSON() ([]byte, error) {
	type iResponse IResponse
	resp := &struct {
		iResponse
		ErrorResp string `json:"error,omitempty"`
	}{
		iResponse: (iResponse)(i),
	}
	if i.Error != nil {
		resp.ErrorResp = i.Error.Error()
	}
	return json.Marshal(resp)
}
