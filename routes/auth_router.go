package routes

import (
	"NeraJima/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(group fiber.Router) {
	router := group.Group("/auth")

	router.Get("/register", auth.Register)
}
