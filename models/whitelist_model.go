package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Database Indices: nil

type Whitelist struct {
	Id        primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OwnerId   primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
	AllowedId primitive.ObjectID `json:"allowedId" bson:"allowedId,omitempty"`
}
