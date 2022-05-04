package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DBClient *mongo.Client

func ConnectDatabase() {
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoUri()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}

	DBClient = client
	fmt.Println("Database connection established...")
}

func GetCollection(databaseName, collectionName string) *mongo.Collection {
	collection := DBClient.Database(databaseName).Collection(collectionName)
	return collection
}
