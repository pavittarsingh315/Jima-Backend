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
	var user models.User
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
	userErr := configs.UserCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if userErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errMessage}})
	}

	profileErr := configs.ProfileCollection.FindOne(ctx, bson.M{"userId": user.Id}).Decode(&profile)
	if profileErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errMessage}})
	}

	c.Locals("user", user)
	c.Locals("profile", profile)

	return c.Next()
}
