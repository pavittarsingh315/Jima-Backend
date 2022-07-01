package users

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/responses"
	"context"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddUserToWhitelist(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistObj models.Whitelist
	var toBeAddedProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == profileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot whitelist yourself."}})
	}

	err := configs.WhitelistCollection.FindOne(ctx, bson.M{"ownerId": reqProfile.Id, "allowedId": profileId}).Decode(&whitelistObj)
	if err == nil { // no error => user is already whitelisted
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User is already whitelisted."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": profileId}).Decode(&toBeAddedProfile)
	if err != nil { // error => user doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Could not whitelist user."}})
	}

	newWhitelistObj := models.Whitelist{
		Id:        primitive.NewObjectID(),
		OwnerId:   reqProfile.Id,
		AllowedId: profileId,
	}

	_, err = configs.WhitelistCollection.InsertOne(ctx, newWhitelistObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	update := bson.M{"numWhitelisted": reqProfile.NumWhitelisted + 1, "lastUpdate": time.Now()}
	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"userId": reqProfile.UserId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "User whitelisted."}})
}

func RemoveUserFromWhitelist(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistObj models.Whitelist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == profileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot blacklist yourself."}})
	}

	err := configs.WhitelistCollection.FindOneAndDelete(ctx, bson.M{"ownerId": reqProfile.Id, "allowedId": profileId}).Decode(&whitelistObj)
	if err != nil { // error => user is not whitelisted
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User is already blacklisted."}})
	}

	update := bson.M{"numWhitelisted": reqProfile.NumWhitelisted - 1, "lastUpdate": time.Now()}
	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"userId": reqProfile.UserId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "User blacklisted."}})
}

func GetWhitelist(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	// search := c.Query("search", "")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "ownerId", Value: reqProfile.Id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "profiles"},
		{Key: "localField", Value: "allowedId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "profile"},
	}}}
	// search stage here
	skipStage := bson.D{{Key: "$skip", Value: (page - 1) * limit}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}}

	aggPipeline := mongo.Pipeline{matchStage, lookupStage, skipStage, limitStage, unwindStage}
	cursor, err := configs.WhitelistCollection.Aggregate(ctx, aggPipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	defer cursor.Close(ctx)

	var whitelistedUsers = []models.MiniProfile{}
	var totalObjects int = 0
	for cursor.Next(ctx) {
		var object struct {
			Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
			OwnerId   primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
			AllowedId primitive.ObjectID `json:"allowedId" bson:"allowedId,omitempty"`
			Profile   models.Profile     `json:"profile" bson:"profile,omitempty"`
		}
		if err := cursor.Decode(&object); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
		}
		var whitelistedUser = models.MiniProfile{
			Id:                 object.Profile.Id,
			Username:           object.Profile.Username,
			Name:               object.Profile.Name,
			MiniProfilePicture: object.Profile.MiniProfilePicture,
		}
		whitelistedUsers = append(whitelistedUsers, whitelistedUser)
		totalObjects++
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    math.Ceil(float64(totalObjects) / float64(limit)),
				"data":         whitelistedUsers,
			},
		},
	)
}
