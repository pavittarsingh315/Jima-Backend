package models

import "time"

type TemporaryObject struct {
	VerificationCode int       `json:"code,omitempty" bson:"code,omitempty"`
	Contact          string    `json:"contact,omitempty" bson:"contact,omitempty"`
	CreatedAt        time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
}
