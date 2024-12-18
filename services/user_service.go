package services

import (
	// "strings"
	// "fmt"
	// "log"
	// "time"
	"encoding/base64"
	"os"
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

func UnpackUser(c *gin.Context, user_id primitive.ObjectID) (string, string, error) {
	// Check if user_id is valid
	var author models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"_id": user_id}).Decode(&author)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return "", "", err
	}
	// The "Avater" field in target_user is a path to the avatar file, not the image data itself
	avatar := author.Avatar
	if avatar == "" {
		avatar = "/tmp/avatar.png"
	}
	picture, err := os.ReadFile(avatar)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, "Failed to read avatar file", err)
		return "", "", err
	}
	// Encode the image data as base64
	encoded := base64.StdEncoding.EncodeToString(picture)
	return encoded, author.Name, nil
}

func UnpackDynamics(c *gin.Context, dynamics []models.Dynamic) []models.DynamicContent {
	var dynamic_contents []models.DynamicContent
	for _, dynamic := range dynamics {
		// Find the content of the dynamic using the dynamic type and the corresponding ID
		var content models.DynamicContent
		content.ID = dynamic.ID
		content.Author = dynamic.Author
		content.Avatar, content.Name, _ = UnpackUser(c, dynamic.Author)
		// Process dynamic content of different types
		if dynamic.Type == 0 { // Classic masterpiece
			// var classic models.Ci
			err := config.MongoClient.Database("serenesong").Collection("Ci").FindOne(c, bson.M{"_id": dynamic.CiId}).Decode(&content.Classic)
			if err != nil {
				utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
				return nil
			}
			// Pack the content and the comments
			content.Type = 0
		} else if dynamic.Type == 1 { // Modern works
			// var modern models.ModernWork
			err := config.MongoClient.Database("serenesong").Collection("UserWorks").FindOne(c, bson.M{"_id": dynamic.UserWorkId}).Decode(&content.Modern)
			if err != nil {
				utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
				return nil
			}
			// Pack the content and the comments
			content.Type = 1
		} else if dynamic.Type == 2 { // Collections
			// find the collection using the collection ID
			result := config.MongoClient.Database("serenesong").Collection("collections").FindOne(c, bson.M{"_id": dynamic.CollectionId})
			if result.Err() != nil {
				utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, result.Err())
				return nil
			}
			// Decode the collection
			var collection models.Collection
			err := result.Decode(&collection)
			if err != nil {
				utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
				return nil
			}
			// Pack the content and the comments
			for _, item := range collection.CollectionItems {
				if item.CiId == dynamic.CollectionCiId {
					content.Comment = item.Comment
					err := config.MongoClient.Database("serenesong").Collection("Ci").FindOne(c, bson.M{"_id": item.CiId}).Decode(&(content.CollectionCi))
					if err != nil {
						utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
						return nil
					}
					break
				}
			}
			// Pack the content and the comments
			content.Type = 2
		} else {
			utils.HandleError(c, http.StatusBadRequest, "Invalid dynamic type", nil)
			return nil
		}
		// Fetch the comments
		for _, comment_id := range dynamic.Comments {
			var comment models.Comment
			err := config.MongoClient.Database("serenesong").Collection("Comments").FindOne(c, bson.M{"_id": comment_id}).Decode(&comment)
			if err != nil {
				utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
				return nil
			}
			// Pack the comment
			packet := models.CommentPacket{
				Comment:   comment,
				CommentId: comment_id,
			}
			packet.Avatar, packet.Name, _ = UnpackUser(c, comment.Commenter)
			content.Comments = append(content.Comments, packet)
		}
		// Fetch the likes
		content.Likes = dynamic.Likes
		// Append the content to the dynamic_contents array
		dynamic_contents = append(dynamic_contents, content)
	}
	// Return the dynamic_contents array
	return dynamic_contents
}

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
	// Unpack the dynamic IDs and fetch the corresponding content
	var dynamic_indices []models.Dynamic
	for _, dynamic_id := range target_user.Dynamics {
		var dynamic models.Dynamic
		err := config.MongoClient.Database("serenesong").Collection("Dynamics").FindOne(c, bson.M{"_id": dynamic_id}).Decode(&dynamic)
		if err != nil {
			utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
			return
		}
		dynamic_indices = append(dynamic_indices, dynamic)
	}
	var dynamics = UnpackDynamics(c, dynamic_indices)
	if dynamics != nil {
		// Return dynamics of the target user
		c.JSON(http.StatusOK, gin.H{"dynamics": dynamics})
	}
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
			utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
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

func ReturnUserInfo(c *gin.Context, user_id string, token string) {
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
	// The "Avater" field in target_user is a path to the avatar file, not the image data itself
	avatar := target_user.Avatar
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
	// Return the base64 encoded image data
	c.JSON(http.StatusOK, gin.H{
		"avatar":    encoded,
		"name":      target_user.Name,
		"signature": target_user.Signature,
	})
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

func SaveNameAvatar(c *gin.Context, token string, name string, avatar string, signature string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Update user name and avatar
	// Make sure the avatar directory exists
	folder := "/tmp/TsingpingYue/avatars"
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
	path := filepath.Join(folder, user.ID.Hex()+".png")
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
				"name":      name,
				"avatar":    path,
				"signature": signature,
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

func ReturnPersonalID(c *gin.Context, token string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Return personal ID of the user
	c.JSON(http.StatusOK, gin.H{
		"personal_id": user.ID,
	})
}
