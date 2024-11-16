package models

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ModernWork struct {
	ID    		primitive.ObjectID	`bson:"_id,omitempty" json:"_id,omitempty"`
	Author 		primitive.ObjectID 	`bson:"author" json:"author"`
	Title 		string 				`bson:"title, omitempty" json:"title,omitempty"`
	// Age         string 				`bson:"age, omitempty" json:"age,omitempty"`
	Content 	[]string 			`bson:"content" json:"content"`
	Cipai   	[]string 			`bson:"cipai" json:"cipai"`
	Xiaoxu 		string 				`bson:"prologue,omitempty" json:"prologue,omitempty"`
	// IsModern 	bool 				`bson:"is_modern" json:"is_modern"`
	IsPublic 	bool 				`bson:"is_public" json:"public"`
	Tags 		[]string 			`bson:"tags,omitempty" json:"tags,omitempty"`
	CreatedAt 	time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt 	time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}