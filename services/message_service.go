package services

import (
	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"time"
)

func GetMessagesIGetHandler(c *gin.Context, token string) {
	// Get user from token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}

	// Get messages
	var messages []models.Message
	cursor, err := config.MongoClient.Database("serenesong").Collection("messages").Find(c, bson.M{"receiver": user.ID})
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		var message models.Message
		err := cursor.Decode(&message)
		if err != nil {
			utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
			return
		}
		messages = append(messages, message)
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func GetMessagesISendHandler(c *gin.Context, token string) {
	// Get user from token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}

	// Get messages
	var messages []models.Message
	cursor, err := config.MongoClient.Database("serenesong").Collection("messages").Find(c, bson.M{"sender": user.ID})
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		var message models.Message
		err := cursor.Decode(&message)
		if err != nil {
			utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
			return
		}
		messages = append(messages, message)
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{"messages": messages})
}

func SendMessageHandler(c *gin.Context, token string, content string, receiver string, replyToMessageId string) {
	// Get user from token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}

	var message models.Message
	message.Sender = user.ID
	message.SenderName = user.Name
	message.Content = content
	message.Time = time.Now()

	// receiver
	receiverId, err := primitive.ObjectIDFromHex(receiver)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidObjID, err)
		return
	}
	var receiverUser models.User
	err = config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": receiverId}).Decode(&receiverUser)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}
	message.Receiver = receiverUser.ID
	message.ReceiverName = receiverUser.Name

	if replyToMessageId != "" {
		replyToMessageIdObj, err := primitive.ObjectIDFromHex(replyToMessageId)
		if err != nil {
			utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidObjID, err)
			return
		}
		message.ReplyTo = replyToMessageIdObj
	}

	// Insert message
	result, err := config.MongoClient.Database("serenesong").Collection("messages").InsertOne(c, message)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoInsert, err)
		return
	}

	// Send response
	c.JSON(http.StatusOK, gin.H{"messageId": result.InsertedID})
}
