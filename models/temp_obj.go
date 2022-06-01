package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Database Indices:
    field {"contact": "1"}, option: {unique: true}
    field {"createdAt": "1"}, option {expireAfterSeconds: 300}
*/

type TemporaryObject struct {
	Id               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	VerificationCode int                `json:"code,omitempty" bson:"code,omitempty"`
	Contact          string             `json:"contact,omitempty" bson:"contact,omitempty"`
	CreatedAt        time.Time          `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
