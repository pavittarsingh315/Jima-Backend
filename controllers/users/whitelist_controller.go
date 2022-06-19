package users

import (
	"NeraJima/responses"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func AddUserToWhitelist(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Whitelisted %s", profileId)}})
}

func RemoveUserFromWhitelist(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Blacklisted %s", profileId)}})
}

func GetWhitelist(c *fiber.Ctx) error {
	// get whitelist for user making request
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Got whitelist"}})
}
