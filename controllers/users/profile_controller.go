package users

import (
	"NeraJima/configs"
	"NeraJima/responses"
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchForUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username, _ := url.QueryUnescape(c.Params("query"))
	regexPattern := fmt.Sprintf("^%s.*", username)

	filter := bson.D{{Key: "username", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Options: "i", Pattern: regexPattern}}}}}
	numDocsRetrievedLimit := options.Find().SetLimit(10)
	fields := options.Find().SetProjection(bson.D{{Key: "username", Value: 1}, {Key: "name", Value: 1}, {Key: "miniProfilePicture", Value: 1}})
	cursor, err := configs.ProfileCollection.Find(ctx, filter, numDocsRetrievedLimit, fields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	var results []struct {
		Id                 primitive.ObjectID `json:"profileId" bson:"_id,omitempty"`
		Username           string             `json:"username,omitempty"`
		Name               string             `json:"name,omitempty"`
		MiniProfilePicture string             `json:"miniProfilePicture,omitempty"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": results}})
}
