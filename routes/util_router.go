package routes

import (
	"NeraJima/controllers/utils"
	"NeraJima/middleware"

	"github.com/gofiber/fiber/v2"
)

func UtilRouter(group fiber.Router) {
	router := group.Group("/util")

	router.Get("/getPresignUrl/profilePicture", middleware.UserAuthHandler, utils.GetProfilePictureUploadUrl)
}
