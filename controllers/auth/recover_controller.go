package auth

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/requests"
	"NeraJima/responses"
	"NeraJima/utils"
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/rivo/uniseg"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RequestPasswordReset(c *fiber.Ctx) error {
	var body requests.RecoveryRequest
	var tempObj models.TemporaryObject
	var user models.User
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	if body.Contact == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", "")) // remove all whitespace and make lowercase

	contactErr := configs.UserCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&user)
	if contactErr != nil { // error => no user with this contact
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	tempObjError := configs.TempObjCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&tempObj)
	if tempObjError == nil { // no error => tempObj with contact exists
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Try again in a few minutes."}})
	}

	code := utils.EncodeToInt(6)
	newTempObj := models.TemporaryObject{
		Id:               primitive.NewObjectID(),
		VerificationCode: code,
		Contact:          body.Contact,
		CreatedAt:        time.Now(),
	}
	_, err := configs.TempObjCollection.InsertOne(ctx, newTempObj)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	contactIsEmail := utils.ValidateEmail(body.Contact)
	if contactIsEmail {
		utils.SendPasswordResetEmail(user.Name, user.Contact, code)
		return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "An email has been sent with a verification code."}})
	} else {
		utils.SendPasswordResetText(code, user.Contact)
		return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "A text has been sent with a verification code."}})
	}
}

func ConfirmResetCode(c *fiber.Ctx) error {
	var body requests.RecoveryRequest
	var tempObj models.TemporaryObject
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	if body.Code == "" || body.Contact == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", "")) // remove all whitespace and make lowercase

	tempObjError := configs.TempObjCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&tempObj)
	if tempObjError != nil { // error => no tempObj was found
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Code has expired. Please restart the reset process."}})
	}

	code, _ := strconv.Atoi(body.Code)
	if tempObj.VerificationCode != code {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Incorrect Code."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Code confirmed."}})
}

func ConfirmPasswordReset(c *fiber.Ctx) error {
	var body requests.RecoveryRequest
	var user models.User
	var profile models.Profile
	var tempObj models.TemporaryObject
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Bad request..."}})
	}

	if body.Code == "" || body.Contact == "" || body.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	passwordLength := uniseg.GraphemeClusterCount(body.Password)
	if passwordLength < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Password too short."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", "")) // remove all whitespace and make lowercase

	tempObjError := configs.TempObjCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&tempObj)
	if tempObjError != nil { // error => no tempObj was found
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Reset time has expired. Please restart the reset process."}})
	}

	code, _ := strconv.Atoi(body.Code)
	if tempObj.VerificationCode != code {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Incorrect Code."}})
	}

	contactErr := configs.UserCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&user)
	if contactErr != nil { // error => no user with this contact
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	if utils.VerifyPassword(user.Password, body.Password) { // old and new passwords match
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "This is your old password...💀"}})
	}

	update := bson.M{"password": utils.HashPassword(body.Password), "lastUpdate": time.Now()}
	_, err := configs.UserCollection.UpdateOne(ctx, bson.M{"_id": user.Id}, bson.M{"$set": update})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	profileErr := configs.ProfileCollection.FindOne(ctx, bson.M{"userId": user.Id}).Decode(&profile)
	if profileErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Account not found."}})
	}

	_, tempObjDelError := configs.TempObjCollection.DeleteOne(ctx, bson.M{"contact": body.Contact})
	if tempObjDelError != nil {
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
