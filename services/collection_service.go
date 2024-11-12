package services

import (
	"Serenesongserver/config"
	"Serenesongserver/models"
	"fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateCollectionHandler(c *gin.Context, collectionName string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	// 创建新的收藏夹
	collection := models.Collection{
		Name:            collectionName,
		CollectionItems: []models.CollectionItem{},
	}

	// 插入数据库并处理错误
	insertResult, err := config.MongoClient.Database("serenesong").Collection("collections").InsertOne(c, collection)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to create collection"})
		return
	}

	// 更新用户的收藏夹列表
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(
		c,
		bson.M{"token": token},
		bson.M{"$push": bson.M{"collections": insertResult.InsertedID}},
	)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to update user collections"})
		fmt.Print(err)
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{
		"message":      "Collection created",
		"collectionID": insertResult.InsertedID,
	})
}

func DeleteCollectionHandler(c *gin.Context, collectionID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	// 删除收藏夹
	_id, _ := primitive.ObjectIDFromHex(collectionID)
	_, err = config.MongoClient.Database("serenesong").Collection("collections").DeleteOne(c, bson.M{"_id": _id})
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to delete collection"})
		return
	}

	// 更新用户的收藏夹列表
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(
		c,
		bson.M{"token": token},
		bson.M{"$pull": bson.M{"collections": _id}},
	)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to update user collections"})
		return
	}

	// 返回成功响应
	c.JSON(200, gin.H{"message": "Collection deleted"})
}

func AddToCollectionHandler(c *gin.Context, collectionID string, ciID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	// 把诗词 ID 加入收藏夹
	collection_object_id, _ := primitive.ObjectIDFromHex(collectionID)
	ci_object_id, _ := primitive.ObjectIDFromHex(ciID)
	updateResult, err := config.MongoClient.Database("serenesong").Collection("collections").UpdateOne(
		c,
		bson.M{"_id": collection_object_id},
		bson.M{"$push": bson.M{"collection_items": models.CollectionItem{CiId: ci_object_id, Comment: ""}}},
	)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to update collection", "error": err.Error()})
		return
	}
	if updateResult.ModifiedCount == 0 {
		c.JSON(404, gin.H{"message": "Collection not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Collection item added"})
}

func RemoveFromCollectionHandler(c *gin.Context, collectionID string, ciID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		c.JSON(404, gin.H{"message": "User not found"})
		return
	}

	// 把收藏夹里面的诗词 ID 踢出去
	collection_object_id, _ := primitive.ObjectIDFromHex(collectionID)
	ci_object_id, _ := primitive.ObjectIDFromHex(ciID)
	updateResult, err := config.MongoClient.Database("serenesong").Collection("collections").UpdateOne(
		c,
		bson.M{"_id": collection_object_id},
		bson.M{"$pull": bson.M{"collection_items": bson.M{"ci_id": ci_object_id}}},
	)
	if err != nil {
		c.JSON(500, gin.H{"message": "Failed to update collection", "error": err.Error()})
		return
	}
	if updateResult.ModifiedCount == 0 {
		c.JSON(404, gin.H{"message": "Collection item not found"})
		return
	}

	c.JSON(200, gin.H{"message": "Collection item removed"})
}
