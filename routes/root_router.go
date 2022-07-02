package routes

import (
	"NeraJima/responses"

	"github.com/gofiber/fiber/v2"
)

func SetupRouter(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/default", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("🚀🚀🚀🚀 - PSJ 05-03-22 8:29 pm")
	})

	AuthRouter(api)
	UserRouter(api)
	UtilRouter(api)
	DevRouter(api)

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusBadRequest).JSON(
			responses.ErrorResponse{
				Status:  fiber.StatusNotFound,
				Message: "Error",
				Data: &fiber.Map{
					"data": "404 not found.",
				},
			},
		)
	})
}
