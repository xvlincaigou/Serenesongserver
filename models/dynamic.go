package models

import "go.mongodb.org/mongo-driver/bson/primitive"

const (
	DYNAMIC_TYPE_CI                 = 0
	DYNAMIC_TYPE_MODERN_WORK        = 1
	DYNAMIC_TYPE_COLLECTION_COMMENT = 2
)

type Dynamic struct {
	ID     primitive.ObjectID `bson:"_id"`
	Author primitive.ObjectID `bson:"author" json:"author"`
	// type 0: 单纯转发一首古代的词。1：单纯转发一首现代的词（其实这个type也可以用来“公开发布自己的诗词”）。2：转发自己的收藏夹批注。
	Type int `bson:"type"`
	// 在type为0的时候，这个字段是对应的古代词的id。
	CiId primitive.ObjectID `bson:"ci_id"`
	// 在type为1的时候，这个字段是对应的现在作品的id。
	UserWorkId primitive.ObjectID `bson:"user_work_id"`
	// 在type为2的时候，这个字段是对应的收藏夹的id。
	CollectionId primitive.ObjectID `bson:"collection_id"`
	// 在type为2的时候，这个字段是对应的收藏对应的词的id。
	CollectionCiId primitive.ObjectID `bson:"collection_ci_id"`
	// 这个动态下面的评论
	Comments []primitive.ObjectID `bson:"comments"`
	Likes    []primitive.ObjectID `bson:"likes"`
}

type DynamicContent struct {
	ID     primitive.ObjectID `bson:"_id"`
	Author primitive.ObjectID `bson:"author" json:"author"`
	// type 0: 单纯转发一首古代的词。1：单纯转发一首现代的词（其实这个type也可以用来“公开发布自己的诗词”）。2：转发自己的收藏夹批注。
	Type int `bson:"type"`
	// 在type为0的时候，这个字段是对应的古代词的id。
	Classic Ci `bson:"classic"`
	// 在type为1的时候，这个字段是对应的现在作品的id。
	Modern ModernWork `bson:"modern"`
	// 在type为2的时候，这个字段是对应的收藏夹对应的词。
	CollectionCi Ci     `bson:"collection_ci"`
	Comment      string `bson:"comment"`
	// 这个动态下面的评论
	Comments []Comment            `bson:"comments"`
	Likes    []primitive.ObjectID `bson:"likes"`
}
