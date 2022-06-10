package utils

import (
	"NeraJima/configs"
	"NeraJima/responses"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func GetProfilePictureUploadUrl(c *fiber.Ctx) error {
	fileName, fileNameErr := configs.GenerateRandS3FileName(64)
	if fileNameErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	largeUploadUrl, largeUrlErr := configs.GenerateS3UploadUrl("profilePictures/large", fileName)
	miniUploadUrl, miniUrlErr := configs.GenerateS3UploadUrl("profilePictures/mini", fileName)

	if largeUrlErr != nil || miniUrlErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	largeFileUrl := strings.Split(largeUploadUrl, "?")[0]
	miniFileUrl := strings.Split(miniUploadUrl, "?")[0]

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"data": &fiber.Map{
					"largeUploadUrl": largeUploadUrl,
					"miniUploadUrl":  miniUploadUrl,
					"largeFileUrl":   largeFileUrl,
					"miniFileUrl":    miniFileUrl,
				},
			},
		},
	)
}
