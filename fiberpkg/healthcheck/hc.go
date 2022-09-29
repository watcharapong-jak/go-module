package healthcheck

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
)

func RunFiberHealthCheck(app *fiber.App) {
	app.Add(http.MethodGet, "/healthz", func(c *fiber.Ctx) error {
		return c.Status(http.StatusOK).JSON("healthy")
	})

	app.Add(http.MethodGet, "/status", func(c *fiber.Ctx) error {

		return c.Status(http.StatusOK).JSON("healthy")
	})
}
