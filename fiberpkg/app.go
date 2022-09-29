package fiberpkg

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/watcharapong-jak/go-module/fiberpkg/appresponse"
	"github.com/watcharapong-jak/go-module/fiberpkg/config"
	"github.com/watcharapong-jak/go-module/fiberpkg/healthcheck"
	"github.com/watcharapong-jak/go-module/fiberpkg/initfiberpkg"
	"net/http"
	"strings"
)

func NewApp(project string, initStage initfiberpkg.InitialAppStruct) *fiber.App {
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return appresponse.JSONResponse(ctx, http.StatusOK, appresponse.IResponse{
				ErrorCode: config.EM.Internal.InternalServerError,
				Error:     err,
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Join(initStage.AllowOrigin, ","),
		AllowMethods:     initStage.AllowMethod,
		AllowHeaders:     initStage.AllowHeader,
		AllowCredentials: true,
		MaxAge:           initStage.MaxAge,
	}))
	healthcheck.RunFiberHealthCheck(app)

	var routesV1 []initfiberpkg.Route

	var prefixPath string
	switch project {
	default:
		routesV1 = initStage.Routes
		prefixPath = "/" + initStage.Path
	}

	var groupPath string
	switch initStage.State {
	case "sit", "dev":
		groupPath = prefixPath + "/v1"
	case "prod":
		groupPath = prefixPath + "/v1"
	default:
		groupPath = "/v1"
	}

	v1 := app.Group(groupPath)

	//v1.Use(middleware.AcceptLanguage)
	//v1.Use(middleware.Recover)
	//v1.Use(middleware.CacheControl)

	for _, ro := range routesV1 {
		if !ro.Test || initStage.State != "prod" {
			v1.Add(ro.Method, ro.Pattern, append(ro.Middleware, ro.Endpoint)...)
			fmt.Printf("[API] %v\t%v%v\n", ro.Method, groupPath, ro.Pattern)
		}
	}

	return app
}
