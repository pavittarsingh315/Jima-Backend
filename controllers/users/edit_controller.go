package users

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/requests"
	"NeraJima/responses"
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rivo/uniseg"
	"go.mongodb.org/mongo-driver/bson"
)

func EditUsername(c *fiber.Ctx) error {
	var body requests.EditProfileRequest
	var user models.User
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	if body.Username == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Username = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Username), " ", "")) // remove all whitespace and make lowercase

	if body.Username == reqProfile.Username {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This is your current username."}})
	}

	usernameLength := uniseg.GraphemeClusterCount(body.Username)
	if usernameLength < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Username too short."}})
	}
	if usernameLength > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Username too long."}})
	}

	usernameErr := configs.UserCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
	if usernameErr == nil { // no error => user with username exists
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Username taken."}})
	}

	update := bson.M{"username": body.Username, "lastUpdate": time.Now()}

	_, updateUserErr := configs.UserCollection.UpdateOne(ctx, bson.M{"_id": reqProfile.UserId}, bson.M{"$set": update}) // use reqProfile.UserId bc user will be undefined
	if updateUserErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	_, updateProfileErr := configs.ProfileCollection.UpdateOne(ctx, bson.M{"userId": reqProfile.UserId}, bson.M{"$set": update}) // use reqProfile.UserId bc user will be undefined
	if updateProfileErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Username updated."}})
}

func EditName(c *fiber.Ctx) error {
	var body requests.EditProfileRequest
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	body.Name = strings.TrimSpace(body.Name) // remove leading and trailing whitespace

	if body.Name == reqProfile.Name {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This is your current name."}})
	}

	nameLength := uniseg.GraphemeClusterCount(body.Name)
	if nameLength > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Name too long."}})
	}

	update := bson.M{"name": body.Name, "lastUpdate": time.Now()}
	_, updateProfileErr := configs.ProfileCollection.UpdateOne(ctx, bson.M{"userId": reqProfile.UserId}, bson.M{"$set": update})
	if updateProfileErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Name updated."}})
}

func EditBio(c *fiber.Ctx) error {
	var body requests.EditProfileRequest
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	body.Bio = strings.TrimSpace(body.Bio) // remove leading and trailing whitespace

	if body.Bio == reqProfile.Bio {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This is your current bio."}})
	}

	bioLength := uniseg.GraphemeClusterCount(body.Bio)
	if len(strings.Split(body.Bio, "\n")) > 6 { // bio has 6 lines max
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Line limit exceeded."}})
	}
	if bioLength > 151 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bio too long."}})
	}

	update := bson.M{"bio": body.Bio, "lastUpdate": time.Now()}
	_, updateProfileErr := configs.ProfileCollection.UpdateOne(ctx, bson.M{"userId": reqProfile.UserId}, bson.M{"$set": update})
	if updateProfileErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Bio updated."}})
}

func EditBlacklistMessage(c *fiber.Ctx) error {
	var body requests.EditProfileRequest
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	body.BlacklistMessage = strings.TrimSpace(body.BlacklistMessage) // remove leading and trailing whitespace

	if body.BlacklistMessage == "" {
		body.BlacklistMessage = "You do not have permission to view these posts!"
	}

	if body.BlacklistMessage == reqProfile.BlacklistMessage {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This is your current blacklist message."}})
	}

	blacklistMessageLength := uniseg.GraphemeClusterCount(body.BlacklistMessage)
	if len(strings.Split(body.BlacklistMessage, "\n")) > 6 { // bio has 6 lines max
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Line limit exceeded."}})
	}
	if blacklistMessageLength > 151 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Message too long."}})
	}

	update := bson.M{"blacklistMessage": body.BlacklistMessage, "lastUpdate": time.Now()}
	_, updateProfileErr := configs.ProfileCollection.UpdateOne(ctx, bson.M{"userId": reqProfile.UserId}, bson.M{"$set": update})
	if updateProfileErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Message updated."}})
}
