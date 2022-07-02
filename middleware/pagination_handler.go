package middleware

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

// Handles page and limit query parameters and appends c.Locals(...)
func PaginationHandler(c *fiber.Ctx) error {
	var page, limit int64 = 1, 10

	if pageParam := c.Query("page"); pageParam != "" {
		pageNum, _ := strconv.Atoi(pageParam)
		if pageNum >= 1 {
			page = int64(pageNum)
		}
	}

	if limitParam := c.Query("limit"); limitParam != "" {
		limitValue, _ := strconv.Atoi(limitParam)
		if limitValue >= 1 {
			if limitValue > 100 {
				limit = 100
			} else {
				limit = int64(limitValue)
			}
		}

	}

	c.Locals("page", page)
	c.Locals("limit", limit)

	return c.Next()
}
