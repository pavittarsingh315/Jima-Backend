package users

import (
	"NeraJima/responses"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

func FollowAUser(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Followed %s", profileId)}})
}

func UnfollowAUser(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Unfollowed %s", profileId)}})
}

func RemoveAFollower(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Removed %s", profileId)}})
}

func GetProfileFollowers(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Got followers for %s", profileId)}})
}

func GetProfileFollowing(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Got following for %s", profileId)}})
}
