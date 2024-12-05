package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Author struct {
	ID    	primitive.ObjectID	`bson:"_id"  json:"_id"`
	Name	string 				`bson:"name" json:"name"`
	Bio 	string 				`bson:"bio"  json:"bio"`
}