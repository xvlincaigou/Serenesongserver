package services

import (
	"encoding/base64"
	"net/http"
	"os"

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

func encodeDefalutAvatar2Base64(c *gin.Context) (string, error) {
	picture, err := os.ReadFile(utils.DefaultAvatarURL)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to read default avatar file", err)
		return "", err
	}
	// Encode the image data as base64
	encoded := base64.StdEncoding.EncodeToString(picture)
	return encoded, nil
}
