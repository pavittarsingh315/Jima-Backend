package users

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/responses"
	"context"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func SearchForUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	username, _ := url.QueryUnescape(c.Params("query"))
	regexPattern := fmt.Sprintf("^%s.*", username)

	filter := bson.D{{Key: "username", Value: bson.D{{Key: "$regex", Value: primitive.Regex{Options: "i", Pattern: regexPattern}}}}}
	numDocsRetrievedLimit := options.Find().SetLimit(10)
	fields := options.Find().SetProjection(bson.D{{Key: "username", Value: 1}, {Key: "name", Value: 1}, {Key: "miniProfilePicture", Value: 1}})
	cursor, err := configs.ProfileCollection.Find(ctx, filter, numDocsRetrievedLimit, fields)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	var results []struct {
		Id                 primitive.ObjectID `json:"profileId" bson:"_id,omitempty"`
		Username           string             `json:"username,omitempty"`
		Name               string             `json:"name,omitempty"`
		MiniProfilePicture string             `json:"miniProfilePicture,omitempty"`
	}
	if err = cursor.All(ctx, &results); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": results}})
}

func GetSearchHistory(c *fiber.Ctx) error {
	var searches models.Search
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := configs.SearchCollection.FindOne(ctx, bson.M{"profileId": reqProfile.Id}).Decode(&searches)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error getting search history..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": searches.Queries}})
}

func RemoveSearchFromHistory(c *fiber.Ctx) error {
	var searches models.Search
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	index, _ := strconv.Atoi(c.Params("index"))
	if index < 0 || index > 21 {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Search out of range..."}})
	}

	err := configs.SearchCollection.FindOne(ctx, bson.M{"profileId": reqProfile.Id}).Decode(&searches)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error getting search history..."}})
	}

	if index >= len(searches.Queries) {
		return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Search out of range..."}})
	}

	newHistory := append(searches.Queries[:index], searches.Queries[index+1:]...)
	update := bson.M{"queries": newHistory}
	_, updateErr := configs.SearchCollection.UpdateOne(ctx, bson.M{"profileId": reqProfile.Id}, bson.M{"$set": update})
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error updating history..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Search Removed."}})
}

func ClearSearchHistory(c *fiber.Ctx) error {
	var searches models.Search
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := configs.SearchCollection.FindOne(ctx, bson.M{"profileId": reqProfile.Id}).Decode(&searches)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error getting search history..."}})
	}

	update := bson.M{"queries": []string{}}
	_, updateErr := configs.SearchCollection.UpdateOne(ctx, bson.M{"profileId": reqProfile.Id}, bson.M{"$set": update})
	if updateErr != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error updating history..."}})
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Searches Cleared."}})
}

func AddSearchHistory(c *fiber.Ctx) error {
	var searches models.Search
	var reqProfile models.Profile = c.Locals("profile").(models.Profile)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query, _ := url.QueryUnescape(c.Params("query"))

	err := configs.SearchCollection.FindOne(ctx, bson.M{"profileId": reqProfile.Id}).Decode(&searches)
	if err != nil {
		newSearch := models.Search{
			Id:        primitive.NewObjectID(),
			ProfileId: reqProfile.Id,
			Queries:   []string{query},
		}
		_, err := configs.SearchCollection.InsertOne(ctx, newSearch)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(responses.ErrorResponse{Status: fiber.StatusBadRequest, Message: "Error", Data: &fiber.Map{"data": "Error. Please try again."}})
		}
	} else {
		currentHistory := searches.Queries
		containsQuery, indexOfQuery := historyContains(currentHistory, query)
		if containsQuery {
			currentHistory = append(currentHistory[:indexOfQuery], currentHistory[indexOfQuery+1:]...) // remove query from slice
		}

		// prepend query to current history
		currentHistory = append(currentHistory, "") // add empty string as last item in slice
		copy(currentHistory[1:], currentHistory)    // move all slice items once to the right
		currentHistory[0] = query                   // make first item the query

		if len(currentHistory) > 22 {
			currentHistory = currentHistory[:len(currentHistory)-1] // remove last item from slice
		}
		update := bson.M{"queries": currentHistory}
		_, updateErr := configs.SearchCollection.UpdateOne(ctx, bson.M{"profileId": reqProfile.Id}, bson.M{"$set": update})
		if updateErr != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Error updating history..."}})
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses.SuccessResponse{Status: fiber.StatusOK, Message: "Success", Data: &fiber.Map{"data": "Search Added."}})
}

func historyContains(history []string, query string) (bool, int) {
	for i, a := range history {
		if a == query {
			return true, i
		}
	}
	return false, 0
}