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
	var subscribers []models.User
	// Get user by ID and hide sensitive information
	for _, user_id := range user.Subscribers {
		var user models.User
		err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": user_id}).Decode(&user)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
			return
		}
		user.Token = ""
		subscribers = append(subscribers, user)
	}
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
	var subscribed_to []models.User
	// Get user by ID and hide sensitive information
	for _, user_id := range user.SubscribedTo {
		var user models.User
		err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": user_id}).Decode(&user)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
			return
		}
		user.Token = ""
		subscribed_to = append(subscribed_to, user)
	}
	c.JSON(http.StatusOK, gin.H{
		"subscribed_to": subscribed_to,
	})
}

func ReturnPublicWorks(c *gin.Context, user_id string, token string) {
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
	// Get public works of the target user
	var public_works []models.ModernWork
	for _, work_id := range target_user.CiWritten {
		// Find work by ID
		var work models.ModernWork
		err := config.MongoClient.Database("serenesong").Collection("UserWorks").FindOne(c, bson.M{"_id": work_id}).Decode(&work)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
			return
		}
		if work.IsPublic {
			public_works = append(public_works, work)
		}
	}
	// Return public works of the target user
	c.JSON(http.StatusOK, gin.H{"public_works": public_works})
}

func ReturnWorks(c *gin.Context, user_id string, token string) {
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
	// Get public works of the target user
	var public_works []models.ModernWork
	for _, work_id := range target_user.CiWritten {
		// Find work by ID
		var work models.ModernWork
		err := config.MongoClient.Database("serenesong").Collection("UserWorks").FindOne(c, bson.M{"_id": work_id}).Decode(&work)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
			return
		}
		public_works = append(public_works, work)
	}
	// Return public works of the target user
	c.JSON(http.StatusOK, gin.H{"public_works": public_works})
}

func ReturnAvatar(c *gin.Context, user_id string, token string) {
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
	// Return avatar of the target user
	c.JSON(http.StatusOK, gin.H{"avatar": target_user.Avatar})
}

func ChangePrivacy(c *gin.Context, work_id string, token string, is_public bool) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Convert user_id to ObjectID
	id, err := primitive.ObjectIDFromHex(work_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	// Find work by ID
	_, err = config.MongoClient.Database("serenesong").Collection("UserWorks").UpdateOne(c, bson.M{"_id": id}, bson.M{"$set": bson.M{"is_public": is_public}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Privacy changed successfully"})
}