package routes

import (
	"github.com/gofiber/fiber/v2"
)

func AuthRouter(group fiber.Router) {
	router := group.Group("/auth")

	router.Get("/register", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("Register route!!!")
	})
}
