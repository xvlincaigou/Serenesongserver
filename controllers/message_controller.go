package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func GetMessagesIGet(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}

	services.GetMessagesIGetHandler(c, token)
}

func GetMessagesISend(c *gin.Context) {
	token := c.Query("token")

	if token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}

	services.GetMessagesISendHandler(c, token)
}

func SendMessage(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}

	token, token_ok := json_data["token"].(string)
	content, content_ok := json_data["content"].(string)
	receiver, receiver_ok := json_data["receiver"].(string)
	replyToMessageId, replyToMessageId_ok := json_data["replyToMessageId"].(string)

	if !token_ok || !content_ok || !receiver_ok || !replyToMessageId_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}

	services.SendMessageHandler(c, token, content, receiver, replyToMessageId)
}

func SubscribeOthers(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}

	token, token_ok := json_data["token"].(string)
	receiver, receiver_ok := json_data["receiver"].(string)
	subscribeOrCancel, subscribeOrCancel_ok := json_data["subscribeOrCancel"].(bool)

	if !token_ok || !receiver_ok || !subscribeOrCancel_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}

	services.SubscribeOthersHandler(c, token, receiver, subscribeOrCancel)
}
