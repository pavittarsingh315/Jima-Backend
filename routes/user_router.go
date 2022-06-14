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
	router.Get("/profile/search/history/get", middleware.UserAuthHandler, users.GetSearchHistory)
	router.Put("/profile/search/history/add/:query", middleware.UserAuthHandler, users.AddSearchHistory)
	router.Put("/profile/search/history/remove/:index", middleware.UserAuthHandler, users.RemoveSearchFromHistory)
	router.Put("/profile/search/history/clear", middleware.UserAuthHandler, users.ClearSearchHistory)
}
