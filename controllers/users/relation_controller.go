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

func FollowAUser(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var relationObject models.Relation
	var toBeFollowedProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	toBeFollowedProfileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == toBeFollowedProfileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot follow yourself."}})
	}

	err := configs.RelationCollection.FindOne(ctx, bson.M{"followedId": toBeFollowedProfileId, "followerId": reqProfile.Id}).Decode(&relationObject)
	if err == nil { // no error => user is already followed
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "You're already following this user."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": toBeFollowedProfileId}).Decode(&toBeFollowedProfile)
	if err != nil { // error => user does not exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User does not exist."}})
	}

	newRelationObj := models.Relation{
		Id:         primitive.NewObjectID(),
		FollowerId: reqProfile.Id,
		FollowedId: toBeFollowedProfileId,
	}

	_, err = configs.RelationCollection.InsertOne(ctx, newRelationObj)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	updateSelf := bson.M{"numFollowing": reqProfile.NumFollowing + 1, "lastUpdate": time.Now()}
	updateFollowed := bson.M{"numFollowers": toBeFollowedProfile.NumFollowers + 1, "lastUpdate": time.Now()}

	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": reqProfile.Id}, bson.M{"$set": updateSelf})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": toBeFollowedProfile.Id}, bson.M{"$set": updateFollowed})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "User followed."}})
}

func UnfollowAUser(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var relationObject models.Relation
	var toBeUnfollowedProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	toBeUnfollowedProfileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == toBeUnfollowedProfileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot unfollow yourself."}})
	}

	err := configs.RelationCollection.FindOneAndDelete(ctx, bson.M{"followedId": toBeUnfollowedProfileId, "followerId": reqProfile.Id}).Decode(&relationObject)
	if err != nil { // error => user is already not followed
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User not followed."}})
	}

	updateSelf := bson.M{"numFollowing": reqProfile.NumFollowing - 1, "lastUpdate": time.Now()}
	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": reqProfile.Id}, bson.M{"$set": updateSelf})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": toBeUnfollowedProfileId}).Decode(&toBeUnfollowedProfile)
	if err != nil { // error => user does not exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User does not exist."}})
	}

	updateUnfollowed := bson.M{"numFollowers": toBeUnfollowedProfile.NumFollowers - 1, "lastUpdate": time.Now()}
	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": toBeUnfollowedProfileId}, bson.M{"$set": updateUnfollowed})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "User unfollowed."}})
}

func RemoveAFollower(c *fiber.Ctx) error {
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	var relationObject models.Relation
	var toBeRemovedProfile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	toBeRemovedProfileId, _ := primitive.ObjectIDFromHex(c.Params("profileId"))

	if reqProfile.Id == toBeRemovedProfileId {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Cannot remove yourself."}})
	}

	err := configs.RelationCollection.FindOneAndDelete(ctx, bson.M{"followedId": reqProfile.Id, "followerId": toBeRemovedProfileId}).Decode(&relationObject)
	if err != nil { // error => user is already not following you
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User is not following you."}})
	}

	updateSelf := bson.M{"numFollowers": reqProfile.NumFollowers - 1, "lastUpdate": time.Now()}
	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": reqProfile.Id}, bson.M{"$set": updateSelf})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	err = configs.ProfileCollection.FindOne(ctx, bson.M{"_id": toBeRemovedProfileId}).Decode(&toBeRemovedProfile)
	if err != nil { // error => user does not exist
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "User does not exist."}})
	}

	updateRemoved := bson.M{"numFollowing": toBeRemovedProfile.NumFollowing - 1, "lastUpdate": time.Now()}
	_, err = configs.ProfileCollection.UpdateOne(ctx, bson.M{"_id": toBeRemovedProfileId}, bson.M{"$set": updateRemoved})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "User removed."}})
}

// TODO: Implement an efficient search/filtering element to the route. Also update the reqProfile's numFollowers if they're incorrect.
func GetProfileFollowers(c *fiber.Ctx) error {
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, err := primitive.ObjectIDFromHex(c.Params("profileId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Invalid id."}})
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "followedId", Value: profileId}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "profiles"},
		{Key: "localField", Value: "followerId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "profile"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "profile.numFollowers", Value: -1}}}}
	skipStage := bson.D{{Key: "$skip", Value: (page - 1) * limit}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "profile._id", Value: 1}, {Key: "profile.username", Value: 1}, {Key: "profile.name", Value: 1}, {Key: "profile.miniProfilePicture", Value: 1}}}}

	aggPipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage, sortStage, skipStage, limitStage, projectStage}
	cursor, err := configs.RelationCollection.Aggregate(ctx, aggPipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer cursor.Close(ctx)

	var followerProfiles = []models.MiniProfile{}
	for cursor.Next(ctx) {
		var object struct {
			Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
			Profile models.MiniProfile `json:"profile" bson:"profile,omitempty"`
		}
		if err := cursor.Decode(&object); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
		}
		followerProfiles = append(followerProfiles, object.Profile)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         followerProfiles,
			},
		},
	)
}

// TODO: Implement an efficient search/filtering element to the route. Also update the reqProfile's numFollowing if they're incorrect.
func GetProfileFollowing(c *fiber.Ctx) error {
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	profileId, err := primitive.ObjectIDFromHex(c.Params("profileId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Invalid id."}})
	}

	matchStage := bson.D{{Key: "$match", Value: bson.D{{Key: "followerId", Value: profileId}}}}
	lookupStage := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "profiles"},
		{Key: "localField", Value: "followedId"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "profile"},
	}}}
	unwindStage := bson.D{{Key: "$unwind", Value: bson.D{{Key: "path", Value: "$profile"}}}}
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "profile.numFollowers", Value: -1}}}}
	skipStage := bson.D{{Key: "$skip", Value: (page - 1) * limit}}
	limitStage := bson.D{{Key: "$limit", Value: limit}}
	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "profile._id", Value: 1}, {Key: "profile.username", Value: 1}, {Key: "profile.name", Value: 1}, {Key: "profile.miniProfilePicture", Value: 1}}}}

	aggPipeline := mongo.Pipeline{matchStage, lookupStage, unwindStage, sortStage, skipStage, limitStage, projectStage}
	cursor, err := configs.RelationCollection.Aggregate(ctx, aggPipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer cursor.Close(ctx)

	var followedProfiles = []models.MiniProfile{}
	for cursor.Next(ctx) {
		var object struct {
			Id      primitive.ObjectID `json:"id" bson:"_id,omitempty"`
			Profile models.MiniProfile `json:"profile" bson:"profile,omitempty"`
		}
		if err := cursor.Decode(&object); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
		}
		followedProfiles = append(followedProfiles, object.Profile)
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    "currently not implemented...", // math.Ceil(float64(totalObjects) / float64(limit))
				"data":         followedProfiles,
			},
		},
	)
}
