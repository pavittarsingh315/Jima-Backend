package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// omitempty means if the field is empty, ignore it.
// i.e, if a field is undefined, it will simply not show the field in the response

type User struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Username    string             `json:"username,omitempty"`
	Password    string             `json:"password,omitempty"`
	Contact     string             `json:"contact,omitempty"`
	Strikes     int                `json:"strikes,omitempty"`
	CreatedDate time.Time          `json:"createdDate,omitempty"`
	LastUpdate  time.Time          `json:"lastUpdate,omitempty"`
	LastLogin   time.Time          `json:"lastLogin,omitempty"`
	BanTill     time.Time          `json:"banTill,omitempty"`
}
