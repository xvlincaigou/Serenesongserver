package services

import (
	// "strings"
	// "fmt"
	// "log"
	// "time"
	"os"
	"encoding/base64"
	"path/filepath"

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
	var dynamics []models.ModernWork
	// := target_user.Dynamics
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
		user.SessionKey = ""
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
		user.SessionKey = ""
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
	// TODO: Return the base64 encoded image data
	// The "Avater" field in target_user is a path to the avatar file, not the image data itself
	picture, err := os.ReadFile(target_user.Avatar)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to read avatar file", err)
		return
	}
	// Encode the image data as base64
	encoded := base64.StdEncoding.EncodeToString(picture)
	// Return the base64 encoded image data
	c.JSON(http.StatusOK, gin.H{"avatar": encoded})
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
	c.JSON(http.StatusOK, gin.H{
		"message": "Privacy changed successfully!",
	})
}

func SaveNameAvatar(c *gin.Context, token string, name string, avatar string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Update user name and avatar
	// Make sure the avatar directory exists
	folder := "/data/TsingpingYue/avatars"
	if err := os.MkdirAll(folder, 0755); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to create avatar directory", err)
		return
	}
	// Decode the base64 encoded image data
	picture, err := base64.StdEncoding.DecodeString(avatar)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid avatar data", err)
		return
	}
	// Save the avatar to the directory
	path := filepath.Join(folder, user.ID.Hex() + ".png")
	err = os.WriteFile(path, picture, 0644)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to save avatar file", err)
		return
	}
	// Update user name and avatar in the database
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(
		c,
		bson.M{"token": token}, 
		bson.M{
			"$set": bson.M{
				"name":   name,
				"avatar": path,
			},
		},
	)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	// Return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "User name and avatar updated successfully!",
	})
}