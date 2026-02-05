package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Reflection struct {
	ID 	  primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserID  primitive.ObjectID `bson:"user_id" json:"user_id"`
	Mood   int             `bson:"mood" json:"mood"`
	Note   string          `bson:"note" json:"note"`
	Date  time.Time       `bson:"date" json:"date"`
	CreatedAt time.Time             `bson:"created_at" json:"created_at"`
}