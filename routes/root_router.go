package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	app.Get("/default", func(c *fiber.Ctx) error {
		return c.SendString("🚀🚀🚀🚀 - PSJ 05-03-22 8:29 pm")
	})

	// api := app.Group("/api")
}
