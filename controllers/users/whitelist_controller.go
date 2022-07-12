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

// TODO: Implement an efficient search/filtering element to the route. Also update the reqProfile's numWhitelisted if they're incorrect.
func GetWhitelist(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "ownerId", Value: reqProfile.Id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "profiles"},
		{Key: "localField", Value: "allowedId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "profile"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "profile.numFollowers", Value: -1}}}}
	skipStage := bson.D{{Key: "$skip", Value: (page - 1) * limit}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "profile._id", Value: 1}, {Key: "profile.username", Value: 1}, {Key: "profile.name", Value: 1}, {Key: "profile.miniProfilePicture", Value: 1}}}}

	aggPipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage, sortStage, skipStage, limitStage, projectStage}
	cursor, err := configs.WhitelistCollection.Aggregate(ctx, aggPipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	defer cursor.Close(ctx)

	var whitelistedUsers = []models.MiniProfile{}
	for cursor.Next(ctx) {
		var object struct {
			Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
			Profile models.MiniProfile `json:"profile" bson:"profile,omitempty"`
		}
		if err := cursor.Decode(&object); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
		}
		whitelistedUsers = append(whitelistedUsers, object.Profile)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         whitelistedUsers,
			},
		},
	)
}

// TODO: Implement an efficient search/filtering element to the route.
func GetWhitelistSubscriptions(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "allowedId", Value: reqProfile.Id}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "profiles"},
		{Key: "localField", Value: "ownerId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "profile"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "profile.username", Value: 1}}}} // sort alphabetically by username
	skipStage := bson.D{{Key: "$skip", Value: (page - 1) * limit}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "profile._id", Value: 1}, {Key: "profile.username", Value: 1}, {Key: "profile.name", Value: 1}, {Key: "profile.miniProfilePicture", Value: 1}}}}

	aggPipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage, sortStage, skipStage, limitStage, projectStage}
	cursor, err := configs.WhitelistCollection.Aggregate(ctx, aggPipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	defer cursor.Close(ctx)

	var subscriptions = []models.MiniProfile{}
	for cursor.Next(ctx) {
		var object struct {
			Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
			Profile models.MiniProfile `json:"profile" bson:"profile,omitempty"`
		}
		if err := cursor.Decode(&object); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
		}
		subscriptions = append(subscriptions, object.Profile)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         subscriptions,
			},
		},
	)
}

func LeaveWhitelist(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistObj models.Whitelist
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == profileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot leave your own whitelist."}})
	}

	err := configs.WhitelistCollection.FindOneAndDelete(ctx, bson.M{"ownerId": profileId, "allowedId": reqProfile.Id}).Decode(&whitelistObj)
	if err != nil { // error => user never had us whitelisted
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This user did not have you whitelisted."}})
	}

	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": profileId}, bson.M{"$inc": bson.M{"numWhitelisted": -1}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Left whitelist successfully."}})
}
