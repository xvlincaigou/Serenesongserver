package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Ci struct {
	ID    		primitive.ObjectID	`bson:"_id" json:"_id"`
	Author 		string 				`bson:"author" json:"author"`
	Title 		string 				`bson:"title, omitempty" json:"title"`
	Age         string 				`bson:"age, omitempty" json:"age"`
	Content 	[]string 			`bson:"content" json:"content"`
	Cipai   	[]string 			`bson:"cipai" json:"cipai"`
	Xiaoxu 		string 				`bson:"prologue" json:"prologue"`
	// IsModern 	bool 				`bson:"is_modern" json:"is_modern"`
	// Public 		bool 				`bson:"public" json:"public"`
	Tags 		[]string 			`bson:"tags" json:"tags"`
	// CreatedAt 	time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
}