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

func InviteToWhitelist(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistObj models.Whitelist
	var whitelistRelationObj models.WhitelistRelation
	var toBeInvitedProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == profileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot invite yourself."}})
	}

	err := configs.WhitelistRelationCollection.FindOne(ctx, bson.M{"senderId": reqProfile.Id, "receiverId": profileId, "type": "Invite"}).Decode(&whitelistRelationObj)
	if err == nil { // no error => invite already sent
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Invite already sent."}})
	}

	err = configs.WhitelistCollection.FindOne(ctx, bson.M{"ownerId": reqProfile.Id, "allowedId": profileId}).Decode(&whitelistObj)
	if err == nil { // no error => user is already whitelisted
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User is already whitelisted."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": profileId}).Decode(&toBeInvitedProfile)
	if err != nil { // error => user doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot invite a user that doesn't exist."}})
	}

	newWhitelistRelationObj, err := models.NewWhitelistRelation(primitive.NewObjectID(), reqProfile.Id, profileId, "Invite")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	_, err = configs.WhitelistRelationCollection.InsertOne(ctx, newWhitelistRelationObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	// TODO: create notification for receiver user

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Invite sent."}})
}

func RevokeWhitelistInvite(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistRelationObj models.WhitelistRelation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inviteId, _ := primitive.ObjectIDFromHex(c.Params("inviteId"))

	err := configs.WhitelistRelationCollection.FindOneAndDelete(ctx, bson.M{"_id": inviteId, "senderId": reqProfile.Id, "type": "Invite"}).Decode(&whitelistRelationObj)
	if err != nil { // error => invite doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User is not invited."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Invite revoked."}})
}

func AcceptWhitelistInvite(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistRelationObj models.WhitelistRelation
	var senderProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inviteId, _ := primitive.ObjectIDFromHex(c.Params("inviteId"))

	err := configs.WhitelistRelationCollection.FindOneAndDelete(ctx, bson.M{"_id": inviteId, "receiverId": reqProfile.Id, "type": "Invite"}).Decode(&whitelistRelationObj)
	if err != nil { // error => invite doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This user has not invited you."}})
	}

	err = configs.ProfileCollection.FindOneAndUpdate(ctx, bson.M{"_id": whitelistRelationObj.SenderId}, bson.M{"$inc": bson.M{"numWhitelisted": 1}}).Decode(&senderProfile)
	if err != nil { // error => sender doesn't exist
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Cannot join a deleted user's whitelist."}})
	}

	newWhitelistObj := models.Whitelist{
		Id:          primitive.NewObjectID(),
		OwnerId:     whitelistRelationObj.SenderId,
		AllowedId:   reqProfile.Id,
		CreatedDate: time.Now(),
	}
	_, err = configs.WhitelistCollection.InsertOne(ctx, newWhitelistObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Invite accepted."}})
}

func DeclineWhitelistInvite(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistRelationObj models.WhitelistRelation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	inviteId, _ := primitive.ObjectIDFromHex(c.Params("inviteId"))

	err := configs.WhitelistRelationCollection.FindOneAndDelete(ctx, bson.M{"_id": inviteId, "receiverId": reqProfile.Id, "type": "Invite"}).Decode(&whitelistRelationObj)
	if err != nil { // error => invite doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This user has not invited you."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Invite declined."}})
}

func RequestWhitelistEntry(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistObj models.Whitelist
	var whitelistRelationObj models.WhitelistRelation
	var toBeRequestedProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == profileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot invite yourself."}})
	}

	err := configs.WhitelistRelationCollection.FindOne(ctx, bson.M{"senderId": reqProfile.Id, "receiverId": profileId, "type": "Request"}).Decode(&whitelistRelationObj)
	if err == nil { // no error => request already sent
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Request already sent."}})
	}

	err = configs.WhitelistCollection.FindOne(ctx, bson.M{"ownerId": profileId, "allowedId": reqProfile.Id}).Decode(&whitelistObj)
	if err == nil { // no error => user already has us whitelisted
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "You're already in this user's whitelist."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": profileId}).Decode(&toBeRequestedProfile)
	if err != nil { // error => user doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot request a user that doesn't exist."}})
	}

	newWhitelistRelationObj, err := models.NewWhitelistRelation(primitive.NewObjectID(), reqProfile.Id, profileId, "Request")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	_, err = configs.WhitelistRelationCollection.InsertOne(ctx, newWhitelistRelationObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	// TODO: create notification for receiver user

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Request sent."}})
}

func CancelWhitelistEntryRequest(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistRelationObj models.WhitelistRelation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestId, _ := primitive.ObjectIDFromHex(c.Params("requestId"))

	err := configs.WhitelistRelationCollection.FindOneAndDelete(ctx, bson.M{"_id": requestId, "senderId": reqProfile.Id, "type": "Request"}).Decode(&whitelistRelationObj)
	if err != nil { // error => request doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "You have not requested to join this user's whitelist."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Request canceled."}})
}

func AcceptWhitelistEntryRequest(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var senderProfile models.Profile
	var whitelistRelationObj models.WhitelistRelation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestId, _ := primitive.ObjectIDFromHex(c.Params("requestId"))

	err := configs.WhitelistRelationCollection.FindOneAndDelete(ctx, bson.M{"_id": requestId, "receiverId": reqProfile.Id, "type": "Request"}).Decode(&whitelistRelationObj)
	if err != nil { // error => request doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This user has not requested to be on your whitelist."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": whitelistRelationObj.SenderId}).Decode(&senderProfile)
	if err != nil { // error => sender doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot admit a user that doesn't exist."}})
	}

	newWhitelistObj := models.Whitelist{
		Id:          primitive.NewObjectID(),
		OwnerId:     reqProfile.Id,
		AllowedId:   whitelistRelationObj.SenderId,
		CreatedDate: time.Now(),
	}
	_, err = configs.WhitelistCollection.InsertOne(ctx, newWhitelistObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": reqProfile.Id}, bson.M{"$inc": bson.M{"numWhitelisted": 1}})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Request accepted."}})
}

func DeclineWhitelistEntryRequest(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var whitelistRelationObj models.WhitelistRelation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	requestId, _ := primitive.ObjectIDFromHex(c.Params("requestId"))

	err := configs.WhitelistRelationCollection.FindOneAndDelete(ctx, bson.M{"_id": requestId, "receiverId": reqProfile.Id, "type": "Request"}).Decode(&whitelistRelationObj)
	if err != nil { // error => request doesn't exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This user has not requested to be on your whitelist."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Request declined."}})
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

func GetWhitelistSentInvites(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": 1,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         "Got Sent Invites",
			},
		},
	)
}

func GetWhitelistReceivedInvites(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": 1,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         "Got Received Invites",
			},
		},
	)
}

func GetWhitelistSentRequests(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": 1,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         "Got Sent Requests",
			},
		},
	)
}

func GetWhitelistReceivedRequests(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": 1,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         "Got Received Requests",
			},
		},
	)
}
