package auth

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/requests"
	"NeraJima/responses"
	"NeraJima/utils"
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Login(c *fiber.Ctx) error {
	var body requests.LoginRequest
	var user models.User
	var profile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	if body.Contact == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", "")) // remove all whitespace and make lowercase

	contactErr := configs.UserCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&user)
	if contactErr != nil { // error => no user with this contact
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	if !utils.VerifyPassword(user.Password, body.Password) { // password doesn't match
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	unixTimeNow := time.Now().Unix()
	unixTimeBan := user.BanTill.Unix()
	if unixTimeNow < unixTimeBan {
		message := fmt.Sprintf("You are banned for %s.", utils.SecondsToString(unixTimeBan-unixTimeNow))
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": message}})
	}

	profileErr := configs.ProfileCollection.FindOne(ctx, bson.M{"userId": user.Id}).Decode(&profile)
	if profileErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	update := bson.M{"lastLogin": time.Now()}
	_, err := configs.UserCollection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	access, refresh := utils.GenAuthTokens(user.Id.Hex())

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"data": &fiber.Map{
					"access":  access,
					"refresh": refresh,
					"profile": profile,
				},
			},
		},
	)
}

func TokenLogin(c *fiber.Ctx) error {
	var body requests.TokenLoginRequest
	var user models.User
	var profile models.Profile
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	if body.AccessToken == "" || body.RefreshToken == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	accessToken, accessBody, accessErr := utils.VerifyAccessToken(body.AccessToken)
	refreshToken, refreshBody, refreshErr := utils.VerifyRefreshToken(body.RefreshToken)
	if accessErr != nil || refreshErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Authentication error..."}})
	}

	if accessBody.UserId != refreshBody.UserId { // token pair are a mismatch
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Authentication error..."}})
	}

	body.AccessToken = accessToken
	body.RefreshToken = refreshToken

	userId, _ := primitive.ObjectIDFromHex(accessBody.UserId)
	userErr := configs.UserCollection.FindOne(ctx, bson.M{"_id": userId}).Decode(&user)
	if userErr != nil { // error => no user with the id from tokens
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	unixTimeNow := time.Now().Unix()
	unixTimeBan := user.BanTill.Unix()
	if unixTimeNow < unixTimeBan {
		message := fmt.Sprintf("You are banned for %s.", utils.SecondsToString(unixTimeBan-unixTimeNow))
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": message}})
	}

	profileErr := configs.ProfileCollection.FindOne(ctx, bson.M{"userId": user.Id}).Decode(&profile)
	if profileErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	update := bson.M{"lastLogin": time.Now()}
	_, err := configs.UserCollection.UpdateOne(ctx, bson.M{"_id": userId}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"data": &fiber.Map{
					"access":  body.AccessToken,
					"refresh": body.RefreshToken,
					"profile": profile,
				},
			},
		},
	)
}
