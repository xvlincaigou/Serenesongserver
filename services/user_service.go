package services

import (
	// "strings"
	// "fmt"
	// "log"
	// "time"

	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"

	// "golang.org/x/text/message/pipeline"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReturnDynamics(c *gin.Context, user_id string, token string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Convert user_id to ObjectID
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	// Check if user_id is valid
	var target_user models.User
	err = config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": id}).Decode(&target_user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}
	// Return dynamics of the target user
	dynamics := target_user.Dynamics
	c.JSON(http.StatusOK, gin.H{
		"dynamics": dynamics,
	})
}

func ReturnCollections(c *gin.Context, user_id string, token string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Convert user_id to ObjectID
	id, err := primitive.ObjectIDFromHex(user_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	// Check if user_id is valid
	var target_user models.User
	err = config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": id}).Decode(&target_user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}
	// Return collections of the target user
	var collections []models.Collection
	for _, collection_id := range target_user.Collections {
		var collection models.Collection
		err = config.MongoClient.Database("serenesong").Collection("collections").FindOne(c, bson.M{"_id": collection_id}).Decode(&collection)
		if err != nil {
			utils.HandleError(c, http.StatusNotFound, "Collection not found", err)
			return
		}
		collections = append(collections, collection)
	}
	c.JSON(http.StatusOK, gin.H{
		"collections": collections,
	})
}

func ReturnSubscribers(c *gin.Context, token string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Return subscribers of the user
	subscribers := user.Subscribers
	c.JSON(http.StatusOK, gin.H{
		"subscribers": subscribers,
	})
}

func ReturnSubscribedTo(c *gin.Context, token string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Return subscribes of the user
	subscribedTo := user.SubscribedTo
	c.JSON(http.StatusOK, gin.H{
		"subscribedTo": subscribedTo,
	})
}