package users

import (
	"NeraJima/configs"
	"NeraJima/models"
	"NeraJima/responses"
	"context"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func SearchForUser(c *fiber.Ctx) error {
	page := c.Locals("page").(int64)
	limit := c.Locals("limit").(int64)
	skip := c.Locals("skip").(int64)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query, _ := url.QueryUnescape(c.Params("query"))
	search := regexp.QuoteMeta(query)

	searchStage := bson.D{{Key: "$match", Value: bson.D{{
		Key: "$or",
		Value: []bson.D{
			{{
				Key:   "username",
				Value: bson.D{{Key: "$regex", Value: primitive.Regex{Options: "i", Pattern: search}}},
			}},
			{{
				Key:   "name",
				Value: bson.D{{Key: "$regex", Value: primitive.Regex{Options: "i", Pattern: search}}},
			}},
		},
	}}}}
	// count # docs here and project value into each doc. then in the for loop, get the # docs and set it equal to totalObjects

	// these 3 stages are optimized: https://stackoverflow.com/questions/24160037/skip-and-limit-in-aggregation-framework
	sortStage := bson.D{{Key: "$sort", Value: bson.D{{Key: "numFollowers", Value: -1}}}}
	limitStage := bson.D{{Key: "$limit", Value: skip + limit}}
	skipStage := bson.D{{Key: "$skip", Value: skip}}

	projectStage := bson.D{{Key: "$project", Value: bson.D{{Key: "username", Value: 1}, {Key: "name", Value: 1}, {Key: "miniProfilePicture", Value: 1}}}}

	aggPipeline := mongo.Pipeline{searchStage, sortStage, limitStage, skipStage, projectStage}
	cursor, err := configs.ProfileCollection.Aggregate(ctx, aggPipeline)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
	}
	defer cursor.Close(ctx)

	var results = []models.MiniProfile{}
	var totalObjects int = 0
	for cursor.Next(ctx) {
		var result models.MiniProfile
		if err := cursor.Decode(&result); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(responses.ErrorResponse{Status: fiber.StatusInternalServerError, Message: "Error", Data: &fiber.Map{"data": "Unexpected error..."}})
		}
		results = append(results, result)
		totalObjects++
	}

	return c.Status(fiber.StatusOK).JSON(
		responses.SuccessResponse{
			Status:  fiber.StatusOK,
			Message: "Success",
			Data: &fiber.Map{
				"current_page": page,
				"last_page":    math.Ceil(float64(totalObjects) / float64(limit)),
				"data":         results,
			},
		},
	)
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
			Id:          primitive.NewObjectID(),
			ProfileId:   reqProfile.Id,
			Queries:     []string{query},
			CreatedDate: time.Now(),
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
