package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// omitempty means if the field is empty, ignore it.
type User struct {
	Id          primitive.ObjectID  `json:"id" bson:"_id,omitempty"`
	Name        string              `json:"name,omitempty" validate:"required"`
	Username    string              `json:"username,omitempty" validate:"required"`
	Password    string              `json:"password,omitempty" validate:"required"`
	Contact     string              `json:"contact,omitempty" validate:"required"`
	Strikes     int                 `json:"strikes,omitempty" validate:"required"`
	CreatedDate time.Time           `json:"createdDate,omitempty" validate:"required"`
	LastUpdate  primitive.Timestamp `json:"lastUpdate,omitempty" validate:"required"`
	LastLogin   primitive.Timestamp `json:"lastLogin,omitempty" validate:"required"`
	BanTill     primitive.Timestamp `json:"banTill,omitempty" validate:"required"`
}
