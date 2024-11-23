package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModernWork struct {
	ID    		primitive.ObjectID	`bson:"_id" json:"_id"`
	Author 		primitive.ObjectID 	`bson:"author" json:"author"`
	Title 		string 				`bson:"title, omitempty" json:"title"`
	// Age         string 				`bson:"age, omitempty" json:"age,omitempty"`
	Content 	[]string 			`bson:"content" json:"content"`
	Cipai   	[]string 			`bson:"cipai" json:"cipai"`
	Xiaoxu 		string 				`bson:"prologue" json:"prologue"`
	// IsModern 	bool 				`bson:"is_modern" json:"is_modern"`
	IsPublic 	bool 				`bson:"is_public" json:"public"`
	Tags 		[]string 			`bson:"tags" json:"tags"`
	CreatedAt 	time.Time           `bson:"created_at" json:"created_at"`
	UpdatedAt 	time.Time           `bson:"updated_at" json:"updated_at"`
}