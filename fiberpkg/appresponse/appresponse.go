package appresponse

import (
	"github.com/gofiber/fiber/v2"
)

// JSONResponse is a function to return response in JSON-format
func JSONResponse(c *fiber.Ctx, status int, v IResponse) error {
	//v.ErrorCode = v.ErrorCode.WithLocale(c)
	return c.Status(status).JSON(v)
}
