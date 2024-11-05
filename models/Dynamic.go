package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dynamic struct {
	// type 0: 诗词评论, 1: 诗词收藏, 2: 诗词发布
	Type     int                `bson:"type"`
    // 是否私密发布
	Secret   bool               `bson:"secret"`
    // 这个动态所关联的诗词ID
	CiId     primitive.ObjectID `bson:"ci_id"`
    // 这个动态下面的评论
	Comments primitive.ObjectID `bson:"comments"`
}
