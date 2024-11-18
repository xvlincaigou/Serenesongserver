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

func CreateCollectionHandler(c *gin.Context, collectionName string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 创建新的收藏夹
	collection := models.Collection{
		ID:              primitive.NewObjectID(),
		Name:            collectionName,
		CollectionItems: []models.CollectionItem{},
	}

	// 插入数据库并处理错误
	insertResult, err := config.MongoClient.Database("serenesong").Collection("collections").InsertOne(c, collection)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoInsert, err)
		return
	}

	// 更新用户的收藏夹列表
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(
		c,
		bson.M{"token": token},
		bson.M{"$push": bson.M{"collections": insertResult.InsertedID}},
	)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"collectionID": insertResult.InsertedID})
}

func DeleteCollectionHandler(c *gin.Context, collectionID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 删除收藏夹
	_id, _ := primitive.ObjectIDFromHex(collectionID)
	_, err = config.MongoClient.Database("serenesong").Collection("collections").DeleteOne(c, bson.M{"_id": _id})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDelete, err)
		return
	}

	// 更新用户的收藏夹列表
	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(
		c,
		bson.M{"token": token},
		bson.M{"$pull": bson.M{"collections": _id}},
	)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}

	// 返回成功响应
	c.JSON(http.StatusOK, gin.H{"message": "Collection deleted"})
}

func AddToCollectionHandler(c *gin.Context, collectionID string, ciID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
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
	if err != nil || updateResult.ModifiedCount == 0 {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection item added"})
}

func RemoveFromCollectionHandler(c *gin.Context, collectionID string, ciID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
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
	if err != nil || updateResult.ModifiedCount == 0 {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection item removed"})
}

func GetAllCollectionsHandler(c *gin.Context, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 查找用户的收藏夹列表
	var collections []models.Collection
	cursor, err := config.MongoClient.Database("serenesong").Collection("collections").Find(c, bson.M{"_id": bson.M{"$in": user.Collections}})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	if err = cursor.All(c, &collections); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"collections": collections})
}

func GetAllCollectionItemsHandler(c *gin.Context, collectionID string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 验证收藏夹是否属于用户
	_id, _ := primitive.ObjectIDFromHex(collectionID)
	collectionOwned := false
	for _, id := range user.Collections {
		if id == _id {
			collectionOwned = true
			break
		}
	}

	if !collectionOwned {
		utils.HandleError(c, http.StatusForbidden, utils.ErrMsgPermission, nil)
		return
	}

	// 查找收藏夹
	var collection models.Collection
	err = config.MongoClient.Database("serenesong").Collection("collections").FindOne(c, bson.M{"_id": _id}).Decode(&collection)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"collectionItems": collection.CollectionItems})
}

func ModifyCollectionCommentHandler(c *gin.Context, ciID string, comment string, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 将 ciID 转换为 ObjectID
	ciObjectID, _ := primitive.ObjectIDFromHex(ciID)
	// 批量查询用户的收藏夹，并使用 $elemMatch 提高效率
	filter := bson.M{
		"_id": bson.M{"$in": user.Collections},
		"collection_items": bson.M{
			"$elemMatch": bson.M{"ciId": ciObjectID},
		},
	}
	update := bson.M{
		"$set": bson.M{"collection_items.$.comment": comment},
	}

	result := config.MongoClient.Database("serenesong").Collection("collections").FindOneAndUpdate(c, filter, update)
	if result.Err() != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, result.Err())
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Collection item comment modified"})
}

func GetCollectionItemCountHandler(c *gin.Context, token string) {
	// 查找用户信息
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 如果用户没有收藏夹，直接返回0
	if len(user.Collections) == 0 {
		c.JSON(http.StatusOK, gin.H{"count": 0})
		return
	}

	// 使用简化的聚合查询
	pipeline := []bson.M{
		{
			"$match": bson.M{
				"_id": bson.M{"$in": user.Collections},
			},
		},
		{
			"$group": bson.M{
				"_id": nil,
				"totalCount": bson.M{
					"$sum": bson.M{"$size": "$collection_items"}, // 修改字段名为 collection_items
				},
			},
		},
	}

	var result []bson.M
	cursor, err := config.MongoClient.Database("serenesong").Collection("collections").Aggregate(c, pipeline)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	defer cursor.Close(c)

	if err = cursor.All(c, &result); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}

	count := 0
	if len(result) > 0 {
		count = int(result[0]["totalCount"].(int32))
	}

	c.JSON(http.StatusOK, gin.H{"count": count})
}
