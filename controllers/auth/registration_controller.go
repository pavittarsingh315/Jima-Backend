package auth

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/requests"
	"NeraJima/responses"
	"NeraJima/utils"
	"context"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func InitiateRegistration(c *fiber.Ctx) error {
	var body requests.RegistrationRequest
	var user models.User
	var tempObj models.TemporaryObject
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	if body.Name == "" || body.Username == "" || body.Password == "" || body.Contact == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", ""))   // remove all whitespace and make lowercase
	body.Username = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Username), " ", "")) // remove all whitespace and make lowercase

	// likelyhood of username being taken is high so we check it first to avoid unnecessary queries/computations
	usernameErr := configs.UserCollection.FindOne(ctx, bson.M{"username": body.Username}).Decode(&user)
	if usernameErr == nil { // no error => user with username exists
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Username taken."}})
	}

	if len(body.Username) < 6 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Username too short."}})
	}
	if len(body.Username) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Username too long."}})
	}
	if len(body.Name) > 30 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Name too long."}})
	}
	if len(body.Contact) > 50 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Contact too long."}})
	}
	if len(body.Password) < 10 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Password too short."}})
	}

	tempObjError := configs.TempObjCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&tempObj)
	if tempObjError == nil { // no error => tempObj with contact exists
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Try again in a few minutes."}})
	}

	contactIsEmail := utils.ValidateEmail(body.Contact)
	contactErr := configs.UserCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&user)
	if contactErr == nil { // no error => user with contact exists
		errorMsg := "Contact already in use."
		if contactIsEmail {
			errorMsg = "Email address already in use."
		}
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": errorMsg}})
	}

	code := utils.EncodeToInt(6)
	newTempObj := models.TemporaryObject{
		VerificationCode: code,
		Contact:          body.Contact,
		CreatedAt:        time.Now(),
	}
	_, err := configs.TempObjCollection.InsertOne(ctx, newTempObj)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
	}

	if contactIsEmail {
		utils.SendRegistrationEmail(body.Name, body.Contact, code)
		return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "An email has been sent with a verification code."}})
	} else {
		utils.SendRegistrationText(code, body.Contact)
		return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "A text has been sent with a verification code."}})
	}
}

func FinalizeRegistration(c *fiber.Ctx) error {
	var body requests.RegistrationRequest
	var tempObj models.TemporaryObject
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	parserErr := c.BodyParser(&body)
	if parserErr != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	if body.Code == 0 || body.Name == "" || body.Username == "" || body.Password == "" || body.Contact == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", ""))   // remove all whitespace and make lowercase
	body.Username = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Username), " ", "")) // remove all whitespace and make lowercase

	tempObjError := configs.TempObjCollection.FindOne(ctx, bson.M{"contact": body.Contact}).Decode(&tempObj)
	if tempObjError != nil { // error => no tempObj was found
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Code has expired. Please restart the registration process."}})
	}

	if tempObj.VerificationCode != body.Code {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Incorrect Code."}})
	}

	newUser := models.User{
		Id:          primitive.NewObjectID(),
		Name:        body.Name,
		Username:    body.Username,
		Password:    utils.HashPassword(body.Password),
		Contact:     body.Contact,
		Strikes:     0,
		CreatedDate: time.Now(),
		LastLogin:   time.Now(),
		LastUpdate:  time.Now(),
		BanTill:     time.Now(),
	}
	_, userErr := configs.UserCollection.InsertOne(ctx, newUser)
	if userErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	// create profile here

	_, err := configs.TempObjCollection.DeleteOne(ctx, bson.M{"contact": body.Contact})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	access, refresh := utils.GenAuthTokens(newUser.Id.Hex())

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"data": &fiber.Map{
					"access":  access,
					"refresh": refresh,
					"user":    newUser,
				},
			},
		},
	)
}
