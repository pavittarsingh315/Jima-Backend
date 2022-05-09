package auth

import (
	"NeraJima/configs"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

var UserCollection *mongo.Collection = configs.GetCollection("Authentication", "User")

func Register(c *fiber.Ctx) error {
	return c.Status(fiber.StatusOK).SendString("Register route!!!")
}
