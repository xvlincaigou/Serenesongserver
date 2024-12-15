package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Message struct {
	Sender       primitive.ObjectID `json:"sender"`
	SenderName   string             `json:"senderName"`
	Receiver     primitive.ObjectID `json:"receiver"`
	ReceiverName string             `json:"receiverName"`
	Time         time.Time          `json:"time"`
	Content      string             `json:"content"`
	ReplyTo      primitive.ObjectID `json:"replyTo"`
}
