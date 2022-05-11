package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
Database Indices:
    field {"contact": "1"}, option: {unique: true}
    field {"username": "1"}, option: {unique: true}
*/

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty" bson:"name,omitempty"`
	Username    string             `json:"username,omitempty" bson:"username,omitempty"`
	Password    string             `json:"password,omitempty" bson:"password,omitempty"`
	Contact     string             `json:"contact,omitempty" bson:"contact,omitempty"`
	Strikes     int                `json:"strikes,omitempty" bson:"strikes,omitempty"`
	CreatedDate time.Time          `json:"createdDate,omitempty" bson:"createdDate,omitempty"`
	LastUpdate  time.Time          `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	LastLogin   time.Time          `json:"lastLogin,omitempty" bson:"lastLogin,omitempty"`
	BanTill     time.Time          `json:"banTill,omitempty" bson:"banTill,omitempty"`
}

// omitempty means if the field is empty, ignore it.
// i.e, if a field is undefined, it will simply not show the field in the response
