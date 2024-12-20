package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	Sender       primitive.ObjectID `bson:"sender" json:"sender"`
	SenderName   string             `bson:"senderName" json:"senderName"`
	Receiver     primitive.ObjectID `bson:"receiver" json:"receiver"`
	ReceiverName string             `bson:"receiverName" json:"receiverName"`
	Time         time.Time          `bson:"time" json:"time"`
	Content      string             `bson:"content" json:"content"`
	ReplyTo      primitive.ObjectID `bson:"replyTo" json:"replyTo"`
}
