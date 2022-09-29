package middleware

import (
	"github.com/gofiber/fiber/v2"
)

func AcceptLanguage() fiber.Handler {
	return func(c *fiber.Ctx) error {
		lang := c.Get("Accept-Language")
		if lang == "" {
			lang = "en"
		}
		c.Locals("locale", lang)
		return c.Next()
	}
}
