package services

import (
	"net/http"

	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetCiByIdHandler 处理获取 Ci 详情的请求
func GetCiByIdHandler(c *gin.Context, _id string) {
	ObjectId, err := primitive.ObjectIDFromHex(_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}

	var ci models.Ci
	err = config.MongoClient.Database("serenesong").Collection("Ci").FindOne(c, bson.M{"_id": ObjectId}).Decode(&ci)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}

	c.JSON(http.StatusOK, ci)
}
