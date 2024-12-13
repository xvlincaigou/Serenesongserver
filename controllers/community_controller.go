package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 发布动态
func PublishDynamic(c *gin.Context) {
	var jsonData bson.M
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}

	token, tokenOk := jsonData["token"].(string)
	_Type, _TypeOk := jsonData["Type"].(float64)
	_id, _idOk := jsonData["_id"].(string)

	if !tokenOk || !_TypeOk || !_idOk {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	Type := int(_Type)
	services.PublishDynamicHandler(c, token, Type, _id)
}
