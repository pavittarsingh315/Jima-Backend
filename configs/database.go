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

	UserCollection = client.Database("Authentication").Collection("users")
	TempObjCollection = client.Database("Authentication").Collection("temporaryobjects")
	ProfileCollection = client.Database("Profiles").Collection("profiles")
	fmt.Println("Database connection established...")
}
