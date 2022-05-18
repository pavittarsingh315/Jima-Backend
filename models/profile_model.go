package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Database Indices:
    field {"username": "1"}, option: {unique: true}
    field {"userId": "1"}, option: {unique: true}
*/

type Profile struct {
	Id                 primitive.ObjectID `json:"profileId" bson:"_id,omitempty"`
	UserId             primitive.ObjectID `json:"userId" bson:"userId,omitempty"`
	Username           string             `json:"username,omitempty" bson:"username,omitempty"`
	Name               string             `json:"name" bson:"name,omitempty"`
	Bio                string             `json:"bio" bson:"bio,omitempty"`
	BlacklistMessage   string             `json:"blacklistMessage,omitempty" bson:"blacklistMessage,omitempty"`
	ProfilePicture     string             `json:"profilePicture,omitempty" bson:"profilePicture,omitempty"`
	MiniProfilePicture string             `json:"miniProfilePicture,omitempty" bson:"miniProfilePicture,omitempty"`
	NumFollowers       int                `json:"numFollowers" bson:"numFollowers,omitempty"`
	NumFollowing       int                `json:"numFollowing" bson:"numFollowing,omitempty"`
	NumWhitelisted     int                `json:"numWhitelisted" bson:"numWhitelisted,omitempty"`
	CreatedDate        time.Time          `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
	LastUpdate         time.Time          `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}

// bson has omitempty because if we update an obj and don't include all fields, the db will replace those fields with empty
// don't omitempty the int fields cause they are allowed to be zero which go considers as empty.
