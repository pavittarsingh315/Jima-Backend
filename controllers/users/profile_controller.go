package users

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/responses"
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAProfile(c *fiber.Ctx) error {
	var profile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))
	fields := options.FindOne().SetProjection(bson.D{{Key: "userId", Value: 0}, {Key: "miniProfilePicture", Value: 0}, {Key: "lastUpdate", Value: 0}})
	err := configs.ProfileCollection.FindOne(ctx, bson.M{"_id": profileId}, fields).Decode(&profile)
	if err != nil { // error => user with profileId doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Could not retrieve user."}})
	}

	response := struct {
		Profile        models.Profile `json:"profile"`
		AreWhitelisted bool           `json:"areWhitelisted"`
		AreFollowing   bool           `json:"areFollowing"`
	}{
		Profile:        profile,
		AreWhitelisted: false, // TODO: replace this dummy value with actual value
		AreFollowing:   false, // TODO: replace this dummy value with actual value
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": response}})
}
