package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type CollectionItem struct {
	// CollectionItem存放了收藏的内容的ID和自己写的批注
	CiId    primitive.ObjectID `bson:"ci_id"`
	Comment string             `bson:"comment"`
}

type Collection struct {
	Name            string           `bson:"name"`
	CollectionItems []CollectionItem `bson:"collection_items"`
}
