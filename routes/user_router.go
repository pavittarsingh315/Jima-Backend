package routes

import (
	"NeraJima/controllers/users"
	"NeraJima/middleware"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(group fiber.Router) {
	router := group.Group("/user") // domain/api/user

	editRouter(router)
	profileRouter(router)
}

func editRouter(group fiber.Router) {
	router := group.Group("/edit") // domain/api/user/edit

	router.Put("/username", middleware.UserAuthHandler, users.EditUsername)
	router.Put("/name", middleware.UserAuthHandler, users.EditName)
	router.Put("/bio", middleware.UserAuthHandler, users.EditBio)
	router.Put("/blacklistmessage", middleware.UserAuthHandler, users.EditBlacklistMessage)
	router.Put("/profilePicture", middleware.UserAuthHandler, users.EditProfilePicture)
}

func profileRouter(group fiber.Router) {
	router := group.Group("/profile") // domain/api/user/profile

	searchRouter(router)
	router.Get("/get/:profileId", middleware.UserAuthHandler, users.GetAProfile)
	relationRouter(router)
	whitelistRouter(router)
}

func searchRouter(group fiber.Router) {
	router := group.Group("/search") // domain/api/user/profile/search

	router.Get("/user/:query", middleware.UserAuthHandler, middleware.PaginationHandler, users.SearchForUser)
	router.Get("/history/get", middleware.UserAuthHandler, users.GetSearchHistory)
	router.Put("/history/add/:query", middleware.UserAuthHandler, users.AddSearchHistory)
	router.Put("/history/remove/:index", middleware.UserAuthHandler, users.RemoveSearchFromHistory)
	router.Put("/history/clear", middleware.UserAuthHandler, users.ClearSearchHistory)
}

func relationRouter(group fiber.Router) {
	router := group // domain/api/user/profile

	router.Post("/follow/:profileId", middleware.UserAuthHandler, users.FollowAUser)
	router.Delete("/unfollow/:profileId", middleware.UserAuthHandler, users.UnfollowAUser)
	router.Delete("/followers/remove/:profileId", middleware.UserAuthHandler, users.RemoveAFollower)
	router.Get("/followers/:profileId", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetProfileFollowers)
	router.Get("/following/:profileId", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetProfileFollowing)
}

func whitelistRouter(group fiber.Router) {
	router := group.Group("/whitelist") // domain/api/user/profile/whitelist

	router.Post("/invite/:profileId", middleware.UserAuthHandler, users.InviteToWhitelist)
	router.Delete("/invite/revoke/:inviteId", middleware.UserAuthHandler, users.RevokeWhitelistInvite)
	router.Post("/invite/accept/:inviteId", middleware.UserAuthHandler, users.AcceptWhitelistInvite)
	router.Delete("/invite/decline/:inviteId", middleware.UserAuthHandler, users.DeclineWhitelistInvite)

	router.Post("/request/:profileId", middleware.UserAuthHandler, users.RequestWhitelistEntry)
	router.Delete("/request/cancel/:requestId", middleware.UserAuthHandler, users.CancelWhitelistEntryRequest)
	router.Post("/request/accept/:requestId", middleware.UserAuthHandler, users.AcceptWhitelistEntryRequest)
	router.Delete("/request/decline/:requestId", middleware.UserAuthHandler, users.DeclineWhitelistEntryRequest)

	router.Delete("/remove/:profileId", middleware.UserAuthHandler, users.RemoveUserFromWhitelist)

	router.Delete("/leave/:profileId", middleware.UserAuthHandler, users.LeaveWhitelist)

	router.Get("/get", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetWhitelist)
	router.Get("/subscriptions/get", middleware.UserAuthHandler, middleware.PaginationHandler, users.GetWhitelistSubscriptions)
}
