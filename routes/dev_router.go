package routes

import (
	"NeraJima/controllers/dev"
	"NeraJima/middleware"

	"github.com/gofiber/fiber/v2"
)

func DevRouter(group fiber.Router) {
	router := group.Group("/dev")

	// router.Post("/genMockData", dev.GenMockData)
	router.Get("/getMockData", middleware.PaginationHandler, dev.GetMockData)
}
