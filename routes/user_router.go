package routes

import (
	"NeraJima/controllers/users"
	"NeraJima/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(group fiber.Router) {
	router := group.Group("/user")

	router.Put("/edit/username", middleware.UserAuthHandler, users.EditUsername)
	router.Put("/edit/name", middleware.UserAuthHandler, users.EditName)
	router.Put("/edit/bio", middleware.UserAuthHandler, users.EditBio)
	router.Put("/edit/blacklistmessage", middleware.UserAuthHandler, users.EditBlacklistMessage)
	router.Put("/edit/profilePicture", middleware.UserAuthHandler, users.EditProfilePicture)

	router.Get("/profile/search/user/:query", middleware.UserAuthHandler, users.SearchForUser)
}
