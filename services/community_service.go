package services

import (
	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PublishDynamicHandler(c *gin.Context, token string, Type int, _id string) {
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	var dynamic models.Dynamic
	dynamic.ID = primitive.NewObjectID()
	dynamic.Author = user.ID
	dynamic.Type = Type
	dynamic.CiId = primitive.ObjectID{}
	dynamic.UserWorkId = primitive.ObjectID{}
	dynamic.CollectionItemId = primitive.ObjectID{}
	dynamic.Comments = []primitive.ObjectID{}
	switch Type {
	case models.DYNAMIC_TYPE_CI:
		dynamic.CiId, err = primitive.ObjectIDFromHex(_id)
	case models.DYNAMIC_TYPE_MODERN_WORK:
		dynamic.UserWorkId, err = primitive.ObjectIDFromHex(_id)
	case models.DYNAMIC_TYPE_COLLECTION_COMMENT:
		dynamic.CollectionItemId, err = primitive.ObjectIDFromHex(_id)
	}
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInvalidObjID, err)
		return
	}

	result, err := config.MongoClient.Database("serenesong").Collection("dynamics").InsertOne(c, dynamic)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoInsert, err)
		return
	}

	insertedID := result.InsertedID
	id, ok := insertedID.(primitive.ObjectID)
	if !ok {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInvalidObjID, err)
		return
	}
	user.Dynamics = append(user.Dynamics, id)
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$set": bson.M{"dynamics": user.Dynamics}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id.Hex()})
}
