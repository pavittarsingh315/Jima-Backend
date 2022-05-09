package models

import "time"

type TemporaryObject struct {
	VerificationCode int       `json:"code,omitempty"`
	Contact          string    `json:"contact,omitempty"`
	CreatedAt        time.Time `json:"createdAt,omitempty" bson:"createdAt"`
}
