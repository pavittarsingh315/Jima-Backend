package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Database Indices:
    field {"ProfileId": "1"}, option: {unique: true}
*/

type Search struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ProfileId primitive.ObjectID `json:"profileId" bson:"profileId,omitempty"`
	Queries   []string           `json:"queries" bson:"queries,omitempty"`
}
