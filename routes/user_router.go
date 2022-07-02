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

	router.Get("/profile/search/user/:query", middleware.UserAuthHandler, middleware.PaginationHandler, users.SearchForUser)
	router.Get("/profile/search/history/get", middleware.UserAuthHandler, users.GetSearchHistory)
	router.Put("/profile/search/history/add/:query", middleware.UserAuthHandler, users.AddSearchHistory)
	router.Put("/profile/search/history/remove/:index", middleware.UserAuthHandler, users.RemoveSearchFromHistory)
	router.Put("/profile/search/history/clear", middleware.UserAuthHandler, users.ClearSearchHistory)

	router.Get("/profile/get/:profileId", middleware.UserAuthHandler, users.GetAProfile)

	router.Post("/profile/follow/:profileId", middleware.UserAuthHandler, users.FollowAUser)
	router.Delete("/profile/unfollow/:profileId", middleware.UserAuthHandler, users.UnfollowAUser)
	router.Delete("/profile/followers/remove/:profileId", middleware.UserAuthHandler, users.RemoveAFollower)
	router.Get("/profile/followers/:profileId", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetProfileFollowers)
	router.Get("/profile/following/:profileId", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetProfileFollowing)

	router.Post("/profile/whitelist/add/:profileId", middleware.UserAuthHandler, users.AddUserToWhitelist)
	router.Delete("/profile/whitelist/remove/:profileId", middleware.UserAuthHandler, users.RemoveUserFromWhitelist)
	router.Get("/profile/whitelist/get", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetWhitelist)
}
