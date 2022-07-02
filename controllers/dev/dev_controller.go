package dev

import (
	"NeraJima/configs"
	"NeraJima/responses"
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type fakeData struct {
	Price int    `json:"price" bson:"price"`
	Name  string `json:"name" bson:"name"`
}

func GenMockData(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var idk []interface{}

	for i := 0; i < 50; i++ {
		newFakeData := fakeData{
			Price: i,
			Name:  fmt.Sprintf("name %d", i),
		}
		idk = append(idk, newFakeData)
	}

	_, err := configs.DevCollection.InsertMany(ctx, idk)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Success"}})
}

func GetMockData(c *fiber.Ctx) error {
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	search := c.Query("search", "")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"$or": []bson.M{
			{
				"name": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}},
			},
			{
				"price": bson.M{"$regex": primitive.Regex{Pattern: search, Options: "i"}},
			},
		},
	}
	options := options.Find()
	options.SetLimit(limit)
	options.SetSkip((page - 1) * limit)

	totalObjects, err := configs.DevCollection.CountDocuments(ctx, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	cursor, err := configs.DevCollection.Find(ctx, filter, options)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer cursor.Close(ctx)

	var items = []fakeData{}
	if err = cursor.All(ctx, &items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    math.Ceil(float64(totalObjects) / float64(limit)),
				"data":         items,
			},
		},
	)
}
