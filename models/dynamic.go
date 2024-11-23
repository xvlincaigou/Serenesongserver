package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Dynamic struct {
	// type 0: 单纯转发一首诗词，可以是古代的或者现代的, 1: 自己的收藏夹里面写了批注的诗词, 2: 公开自己创作的诗词
	Type int `bson:"type"`
	// 这个动态所关联的诗词ID
	CiId primitive.ObjectID `bson:"ci_id"`
	// 这个动态下面的评论
	Comments primitive.ObjectID `bson:"comments"`
}
