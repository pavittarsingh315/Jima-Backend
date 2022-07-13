package configs

import (
	"context"
	"fmt"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var UserCollection *mongo.Collection
var TempObjCollection *mongo.Collection
var ProfileCollection *mongo.Collection
var SearchCollection *mongo.Collection
var RelationCollection *mongo.Collection
var WhitelistCollection *mongo.Collection
var WhitelistRelationCollection *mongo.Collection

func ConnectDatabase() *mongo.Client {
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

	UserCollection = client.Database("Authentication").Collection("users")
	TempObjCollection = client.Database("Authentication").Collection("temporaryobjects")

	ProfileCollection = client.Database("Profiles").Collection("profiles")
	SearchCollection = client.Database("Profiles").Collection("searches")
	RelationCollection = client.Database("Profiles").Collection("relations")
	WhitelistCollection = client.Database("Profiles").Collection("whitelists")
	WhitelistRelationCollection = client.Database("Profiles").Collection("whitelistRelations")

	fmt.Println("Database connection established...")

	return client
}
