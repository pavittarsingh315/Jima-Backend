package models

import (
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Database Indices: nil

type Whitelist struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	OwnerId     primitive.ObjectID `json:"ownerId" bson:"ownerId,omitempty"`
	AllowedId   primitive.ObjectID `json:"allowedId" bson:"allowedId,omitempty"`
	CreatedDate time.Time          `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
}

// Database Indices: nil

type WhitelistRelation struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	SenderId    primitive.ObjectID `json:"senderId" bson:"senderId,omitempty"`
	ReceiverId  primitive.ObjectID `json:"receiverId" bson:"receiverId,omitempty"`
	Type        string             `json:"type" bson:"type,omitempty"`
	CreatedDate time.Time          `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
}

func NewWhitelistRelation(objectId, senderId, receiverId primitive.ObjectID, objectType string) (WhitelistRelation, error) {
	if objectType != "Invite" && objectType != "Request" {
		return WhitelistRelation{}, errors.New("object type not valid")
	}

	return WhitelistRelation{
		Id:          objectId,
		SenderId:    senderId,
		ReceiverId:  receiverId,
		Type:        objectType,
		CreatedDate: time.Now(),
	}, nil
}
