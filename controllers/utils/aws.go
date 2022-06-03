package utils

import (
	"NeraJima/configs"
	"NeraJima/responses"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetProfilePictureUploadUrl(c *fiber.Ctx) error {
	uploadUrl, err := configs.GenerateS3UploadUrl("profilePictures")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	fileUrl := strings.Split(uploadUrl, "?")[0]

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"data": &fiber.Map{
					"uploadUrl": uploadUrl,
					"fileUrl":   fileUrl,
				},
			},
		},
	)
}
