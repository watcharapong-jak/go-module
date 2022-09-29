package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/watcharapong-jak/go-module/fiberpkg/appresponse"
	"github.com/watcharapong-jak/go-module/fiberpkg/config"

	"net/http"
)

func Recover(c *fiber.Ctx) error {
	defer func(ctx *fiber.Ctx) {
		if rec := recover(); rec != nil {
			//err, ok := rec.(error)
			//if !ok {
			//	err = fmt.Errorf("%v", rec)
			//}
			//stack := make([]byte, 4<<10) // 4KB
			//length := runtime.Stack(stack, false)

			appresponse.JSONResponse(ctx, http.StatusOK, appresponse.IResponse{ErrorCode: config.EM.Internal.InternalServerError, Data: "load test panic"})
		}
	}(c)
	return c.Next()
}
