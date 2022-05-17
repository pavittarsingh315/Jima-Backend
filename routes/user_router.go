package routes

import (
	"NeraJima/controllers/users"
	"NeraJima/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(group fiber.Router) {
	router := group.Group("/user")

	router.Put("/edit/username", middleware.UserAuthHandler, users.EditUsername)
}
