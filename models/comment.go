package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Comment 存放了动态的评论，包括评论者、评论内容。这个不是用户对于诗词的批注。
type Comment struct {
	DynamicId primitive.ObjectID `bson:"dynamic_id"`
	Commenter primitive.ObjectID `bson:"commenter"`
	Content   string             `bson:"content"`
}

type CommentPacket struct {
	CommentId primitive.ObjectID `bson:"_id"`
	Name      string             `bson:"name"`
	Avatar    string             `bson:"avatar"`
	Comment   Comment            `bson:"comment"`
}
