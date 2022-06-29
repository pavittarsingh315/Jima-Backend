package users

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/responses"
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func GetProfileFollowers(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Got followers for %s", profileId)}})
}

func GetProfileFollowing(c *fiber.Ctx) error {
	profileId := c.Params("profileId")
	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": fmt.Sprintf("Got following for %s", profileId)}})
}
