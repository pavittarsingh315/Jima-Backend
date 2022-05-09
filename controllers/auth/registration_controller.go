package auth

import (
	"NeraJima/requests"
	"NeraJima/responses"
	"strings"

	"github.com/gofiber/fiber/v2"
)

func InitiateRegistration(c *fiber.Ctx) error {
	var body requests.InitiateRegistrationRequest

	err := c.BodyParser(&body)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	if body.Name == "" || body.Username == "" || body.Password == "" || body.Contact == "" {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Please include all fields."}})
	}

	body.Contact = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Contact), " ", ""))   // remove all whitespace and make lowercase
	body.Username = strings.ToLower(strings.ReplaceAll(strings.TrimSpace(body.Username), " ", "")) // remove all whitespace and make lowercase

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

	// check if user with username exists
	// check if temp obj with contact exists

	return c.Status(fiber.StatusOK).SendString("Register route!!!")
}
