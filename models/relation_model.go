package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Database Indices: nil

type Relation struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	FollowerId primitive.ObjectID `json:"followerId" bson:"followerId,omitempty"`
	FollowedId primitive.ObjectID `json:"followedId" bson:"followedId,omitempty"`
}
