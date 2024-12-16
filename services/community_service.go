package services

import (
	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"

	"math/rand"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func PublishDynamicHandler(c *gin.Context, token string, Type int, _id string, _id2 string) {
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
	dynamic.CollectionId = primitive.ObjectID{}
	dynamic.CollectionCiId = primitive.ObjectID{}
	dynamic.Comments = []primitive.ObjectID{}
	dynamic.Likes = []primitive.ObjectID{}
	switch Type {
	case models.DYNAMIC_TYPE_CI:
		dynamic.CiId, err = primitive.ObjectIDFromHex(_id)
	case models.DYNAMIC_TYPE_MODERN_WORK:
		dynamic.UserWorkId, err = primitive.ObjectIDFromHex(_id)
	case models.DYNAMIC_TYPE_COLLECTION_COMMENT:
		var err1, err2 error
		dynamic.CollectionId, err1 = primitive.ObjectIDFromHex(_id)
		dynamic.CollectionCiId, err2 = primitive.ObjectIDFromHex(_id2)
		if err1 != nil || err2 != nil {
			utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, nil)
			return
		}
	}
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInvalidObjID, err)
		return
	}

	result, err := config.MongoClient.Database("serenesong").Collection("Dynamics").InsertOne(c, dynamic)
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

func DeletePost(c *gin.Context, token string, post_id string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Convert dynamic_id to ObjectID
	dynamic_id, err := primitive.ObjectIDFromHex(post_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	// Check if dynamic_id is valid
	var dynamic models.Dynamic
	err = config.MongoClient.Database("serenesong").Collection("Dynamics").FindOne(c, bson.M{"_id": dynamic_id}).Decode(&dynamic)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	// Check if the user is the author of the post
	if dynamic.Author.Hex() != user.ID.Hex() {
		utils.HandleError(c, http.StatusForbidden, "You are not the author of this post", err)
		return
	}
	// Delete all comments of the post
	for _, comment_id := range dynamic.Comments {
		_, err = config.MongoClient.Database("serenesong").Collection("Comments").DeleteOne(c, bson.M{"_id": comment_id})
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDelete, err)
			return
		}
	}
	// Delete the post from the "Dynamics" table
	_, err = config.MongoClient.Database("serenesong").Collection("Dynamics").DeleteOne(c, bson.M{"_id": dynamic_id})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDelete, err)
		return
	}
	// Remove the post from the user's "Dynamics" list
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"_id": user.ID}, bson.M{"$pull": bson.M{"dynamics": dynamic_id}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Post deleted successfully",
		"dynamic": dynamic_id.Hex(),
	})
}

func ReturnRandomPosts(c *gin.Context, token string, value int) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Define the sort option
	option := options.Find()
	option.SetSort(bson.D{{Key: "_id", Value: -1}}) // 1 for ascending, -1 for descending
	var dynamic_indices []models.Dynamic
	// Fetch all posts from the "Dynamics" table
	cursor, err := config.MongoClient.Database("serenesong").Collection("Dynamics").Find(c, bson.M{}, option)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	if err := cursor.All(c, &dynamic_indices); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
		return
	}
	// If the number of posts is less than the value parameter, return all posts
	if len(dynamic_indices) <= value {
		dynamics := UnpackDynamics(c, dynamic_indices)
		c.JSON(http.StatusOK, gin.H{"dynamics": dynamics})
		return
	}
	// Randomly select value number of posts from the list
	random_indices := rand.Perm(len(dynamic_indices))[:value]
	chosen_indices := make([]models.Dynamic, value)
	for i, index := range random_indices {
		chosen_indices[i] = dynamic_indices[index]
	}
	dynamics := UnpackDynamics(c, chosen_indices)
	if dynamics != nil {
		// Return dynamics of the target user
		c.JSON(http.StatusOK, gin.H{"dynamics": dynamics})
	}
}

func ReturnFollowingPosts(c *gin.Context, user_id string, token string) {
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
	// Sort the posts by date
	dynamic_ids := target_user.Dynamics
	sort.Slice(dynamic_ids, func(i, j int) bool {
		return dynamic_ids[i].Hex() > dynamic_ids[j].Hex()
	})
	// Unpack the dynamic IDs and fetch the corresponding content
	var dynamic_indices []models.Dynamic
	for _, dynamic_id := range dynamic_ids {
		var dynamic models.Dynamic
		err := config.MongoClient.Database("serenesong").Collection("Dynamics").FindOne(c, bson.M{"_id": dynamic_id}).Decode(&dynamic)
		if err != nil {
			utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
			return
		}
		dynamic_indices = append(dynamic_indices, dynamic)
	}
	// Unpack the posts and return
	dynamics := UnpackDynamics(c, dynamic_indices)
	if dynamics != nil {
		// Return dynamics of the target user
		c.JSON(http.StatusOK, gin.H{"dynamics": dynamics})
	}
}

func SendComment(c *gin.Context, token string, content string, post_id string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Check if dynamic_id is valid
	dynamic_id, err := primitive.ObjectIDFromHex(post_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	var dynamic_index models.Dynamic
	err = config.MongoClient.Database("serenesong").Collection("Dynamics").FindOne(c, bson.M{"_id": dynamic_id}).Decode(&dynamic_index)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	// Create a new comment object
	comment := models.Comment{
		DynamicId: dynamic_id,
		Commenter: user.ID,
		Content:   content,
	}
	// Insert the comment into the "Comments" table
	result, err := config.MongoClient.Database("serenesong").Collection("Comments").InsertOne(c, comment)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoInsert, err)
		return
	}
	// Append the comment to the dynamic's comments list
	// dynamic_index.Comments = append(dynamic_index.Comments, result.InsertedID.(primitive.ObjectID))
	// Update the dynamic in the "Dynamics" table
	_, err = config.MongoClient.Database("serenesong").Collection("Dynamics").UpdateOne(c, bson.M{"_id": dynamic_id}, bson.M{"$push": bson.M{"comments": result.InsertedID.(primitive.ObjectID)}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	// Return the new comment
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment added successfully",
		"comment": comment,
	})
}

func DeleteComment(c *gin.Context, token string, comment_id string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Check if comment_id is valid
	cmt_id, err := primitive.ObjectIDFromHex(comment_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	// Check if the comment exists
	var comment models.Comment
	err = config.MongoClient.Database("serenesong").Collection("Comments").FindOne(c, bson.M{"_id": cmt_id}).Decode(&comment)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	// Check if the comment belongs to the user
	if comment.Commenter.Hex() != user.ID.Hex() {
		utils.HandleError(c, http.StatusForbidden, "You are not authorized to delete this comment", err)
		return
	}
	// Delete the comment from the "Comments" table
	_, err = config.MongoClient.Database("serenesong").Collection("Comments").DeleteOne(c, bson.M{"_id": cmt_id})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDelete, err)
		return
	}
	// Remove the comment from the dynamic's comments list
	_, err = config.MongoClient.Database("serenesong").Collection("Dynamics").UpdateOne(c, bson.M{"_id": comment.DynamicId}, bson.M{"$pull": bson.M{"comments": cmt_id}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	// Return success message
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment deleted successfully",
		"comment": comment,
	})
}

func SendLike(c *gin.Context, token string, post_id string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Check if dynamic_id is valid
	dynamic_id, err := primitive.ObjectIDFromHex(post_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	var dynamic_index models.Dynamic
	err = config.MongoClient.Database("serenesong").Collection("Dynamics").FindOne(c, bson.M{"_id": dynamic_id}).Decode(&dynamic_index)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	// Check if the user has already liked the post
	for _, like := range dynamic_index.Likes {
		if like.Hex() == user.ID.Hex() {
			// utils.HandleError(c, http.StatusForbidden, "You have already liked this post", err)
			c.JSON(http.StatusForbidden, gin.H{
				"message":    "You have already liked this post, nothing happened",
				"dynamic_id": dynamic_id.Hex(),
				"user_id":    user.ID.Hex(),
			})
			return
		}
	}
	// Update the dynamic in the "Dynamics" table
	_, err = config.MongoClient.Database("serenesong").Collection("Dynamics").UpdateOne(c, bson.M{"_id": dynamic_id}, bson.M{"$push": bson.M{"likes": user.ID}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	// Return success message
	c.JSON(http.StatusOK, gin.H{
		"message":    "Post liked successfully",
		"dynamic_id": dynamic_id.Hex(),
		"user_id":    user.ID.Hex(),
	})
}

func DeleteLike(c *gin.Context, token string, post_id string) {
	// Verify user token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgInvalidToken, err)
		return
	}
	// Check if dynamic_id is valid
	dynamic_id, err := primitive.ObjectIDFromHex(post_id)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}
	var dynamic_index models.Dynamic
	err = config.MongoClient.Database("serenesong").Collection("Dynamics").FindOne(c, bson.M{"_id": dynamic_id}).Decode(&dynamic_index)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}
	// Check if the user has already liked the post
	for _, like := range dynamic_index.Likes {
		if like.Hex() == user.ID.Hex() {
			// Remove the like from the dynamic's likes list
			_, err = config.MongoClient.Database("serenesong").Collection("Dynamics").UpdateOne(c, bson.M{"_id": dynamic_id}, bson.M{"$pull": bson.M{"likes": like}})
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
				return
			}
			// Return success message
			c.JSON(http.StatusOK, gin.H{
				"message":    "Dismiss like successfully",
				"dynamic_id": dynamic_id.Hex(),
				"user_id":    user.ID.Hex(),
			})
			return
		}
	}
	// If the user has not liked the post, return
	// utils.HandleError(c, http.StatusForbidden, "You have not liked this post", err)
	c.JSON(http.StatusForbidden, gin.H{
		"message":    "You have not liked this post, nothing happened",
		"dynamic_id": dynamic_id.Hex(),
		"user_id":    user.ID.Hex(),
	})
}