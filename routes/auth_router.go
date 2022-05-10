package routes

import (
	"NeraJima/controllers/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(group fiber.Router) {
	router := group.Group("/auth")

	router.Post("/register/initial", auth.InitiateRegistration)
	router.Post("/register/final", auth.FinalizeRegistration)
}
