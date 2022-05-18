package middleware

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/responses"
	"NeraJima/utils"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type header struct {
	Token  string `reqHeader:"token"`
	UserId string `reqHeader:"userId"`
}

func UserAuthHandler(c *fiber.Ctx) error {
	var reqHeader header
	var profile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	errMessage := "Could not authorize action."

	parserErr := c.ReqHeaderParser(&reqHeader)
	if parserErr != nil || reqHeader.Token == "" || reqHeader.UserId == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errMessage}})
	}

	_, accessBody, accessErr := utils.VerifyAccessTokenNoRefresh(reqHeader.Token) // will return err if expired

	if accessErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errMessage}})
	}

	if accessBody.UserId != reqHeader.UserId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errMessage}})
	}

	userId, _ := primitive.ObjectIDFromHex(accessBody.UserId)
	profileErr := configs.ProfileCollection.FindOne(ctx, bson.M{"userId": userId}).Decode(&profile)
	if profileErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errMessage}})
	}

	c.Locals("profile", profile)

	return c.Next()
}
