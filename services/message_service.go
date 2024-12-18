package services

import (
	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"os"
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
	message.ID = primitive.NewObjectID()
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

func SubscribeOthersHandler(c *gin.Context, token string, receiver string, subscribeOrCancel bool) {
	// 获取用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}

	// 将 receiver 字符串转换为 ObjectID
	receiverId, err := primitive.ObjectIDFromHex(receiver)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidObjID, err)
		return
	}
	// 先查询接收者信息
	var receiverUser models.User
	err = config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": receiverId}).Decode(&receiverUser)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 判断是关注还是取消关注
	if subscribeOrCancel == true {
		// 检查是否已经关注
		for _, subscribedUserId := range user.SubscribedTo {
			if subscribedUserId == receiverId {
				utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgDuplicate, nil)
				return
			}
		}

		// 进行关注操作
		// 将当前用户ID添加到目标用户的 subscribers 列表
		updateReceiver := bson.M{"$addToSet": bson.M{"subscribers": user.ID}}
		_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"_id": receiverId}, updateReceiver)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
			return
		}
		// 将目标用户ID添加到当前用户的 subscribed_to 列表
		updateUser := bson.M{"$addToSet": bson.M{"subscribed_to": receiverId}}
		_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, updateUser)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
			return
		}

		// 发送响应
		c.JSON(http.StatusOK, gin.H{"message": "Success"})
		return
	}
	// 取消关注
	// 检查是否已关注该用户
	found := false
	for _, subscribedUserId := range user.SubscribedTo {
		if subscribedUserId == receiverId {
			// 找到该用户，准备取消关注
			found = true
			// 执行取消关注操作
			updateUser := bson.M{"$pull": bson.M{"subscribed_to": receiverId}}
			_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, updateUser)
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
				return
			}

			// 同时从目标用户的 subscribers 列表中移除当前用户ID
			updateReceiver := bson.M{"$pull": bson.M{"subscribers": user.ID}}
			_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"_id": receiverId}, updateReceiver)
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
				return
			}

			// 发送响应
			c.JSON(http.StatusOK, gin.H{"message": "Success"})
			return
		}
	}
	// 如果没有找到该用户
	if !found {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgMongoFind, nil)
		return
	}
}

func SearchUserByNameHandler(c *gin.Context, name string) {
	// 查询数据库中与 name 匹配的用户
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").
		FindOne(c, bson.M{"name": name}).Decode(&user)

	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	avatar := user.Avatar
	if avatar == "" {
		avatar = "/tmp/avatar.png"
	}
	picture, err := os.ReadFile(avatar)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to read avatar file", err)
		return
	}
	// Encode the image data as base64
	encoded := base64.StdEncoding.EncodeToString(picture)

	// 构建返回的JSON对象，包含头像、昵称和ID
	response := gin.H{
		"id":     user.ID.Hex(),
		"name":   user.Name,
		"avatar": encoded,
	}

	// 返回查询到的用户信息
	c.JSON(http.StatusOK, response)
}
