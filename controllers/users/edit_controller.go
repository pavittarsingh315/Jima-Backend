package users

import (
	"NeraJima/responses"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func EditUsername(c *fiber.Ctx) error {
	fmt.Println(c.Locals("user"))
	fmt.Println(c.Locals("profile"))

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Edit Username!"}})
}
