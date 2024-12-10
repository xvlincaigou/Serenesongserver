package services

import (
	// "strings"
	// "fmt"
	"log"
	"time"

	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	// "golang.org/x/text/message/pipeline"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ReturnRhymes(c *gin.Context) {
	// Get MongoDB client
	db := config.MongoClient.Database("serenesong")
	if db == nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoConnect, nil)
		return
	}
	// Fetch collections
	rhymes_collection := db.Collection("Characters")
	rhymes_cursor, err := rhymes_collection.Find(c, bson.M{})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	// Decode rhymes
	var rhymes []bson.M
	if err := rhymes_cursor.All(c, &rhymes); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
		return
	}
	// Fetch pingze
	pingze_collection := db.Collection("CharacterTune")
	pingze_cursor, err := pingze_collection.Find(c, bson.M{})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	// Decode pingze
	var pingze []bson.M
	if err := pingze_cursor.All(c, &pingze); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"rhymes": rhymes,
		"pingze": pingze,
	})
}

func ReturnYunshu(c *gin.Context) {
	db := config.MongoClient.Database("serenesong")
	if db == nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoConnect, nil)
		return
	}
	// Fetch collections
	rhymes_collection := db.Collection("PingshuiYun")
	rhymes_cursor, err := rhymes_collection.Find(c, bson.M{})
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	// Decode rhymes
	var rhymes []bson.M
	if err := rhymes_cursor.All(c, &rhymes); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"rhymes": rhymes})
}

func ReturnFormat(c *gin.Context, cipai_name string, format_num int) {
	// Get MongoDB "CipaiList" collection
	collection := config.MongoClient.Database("serenesong").Collection("CipaiList")
	// Get cipai info
	filter := bson.M{cipai_name: bson.M{"$exists": true}}
	proj := bson.M{"_id": 0, cipai_name: 1}
	cursor, err := collection.Find(c, filter, options.Find().SetProjection(proj))
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	// Get cipai info
	var cipai_fields []bson.M
	if err := cursor.All(c, &cipai_fields); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	// Check if cipai_name is valid
	if cipai_fields == nil {
		utils.HandleError(c, http.StatusNotFound, "Format not found", nil)
		return
	}
	// Extract formats from cipai info
	for _, field := range cipai_fields {
		if value, ok := field[cipai_name]; ok {
			if formats, ok := value.(bson.M)["formats"]; ok {
				if format_array, ok := formats.(bson.A); ok && len(format_array) > 0 {
					// Directly return the first format
					c.JSON(http.StatusOK, gin.H{"format": format_array[format_num]})
					return
				}
			}
		}
	}
	// If no formats field found in any matching document, return error
	c.JSON(http.StatusOK, gin.H{"error": "No formats field found in any matching document"})
}

// func SaveWork(c *gin.Context, work bson.M, token string) {
// 	// Verify user token
// 	var user models.User
// 	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
// 	if err != nil {
// 		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
// 		return
// 	}
// 	// Save work to database
// 	user_work := models.ModernWork{
// 		ID:        primitive.NewObjectID(),
// 		Author:    work["author"].(primitive.ObjectID),
// 		Title:     work["title"].(string),
// 		Content:   utils.ToStringArray(work["content"]),
// 		Cipai:     utils.ToStringArray(work["cipai"]),
// 		Xiaoxu:    work["xiaoxu"].(string),
// 		IsPublic:  work["is_public"].(bool),
// 		Tags:      utils.ToStringArray(work["tags"]),
// 		CreatedAt: work["created_at"].(time.Time),
// 		UpdatedAt: time.Now(),
// 	}
// 	collection := config.MongoClient.Database("serenesong").Collection("UserWorks")
// 	_, err = collection.InsertOne(c, user_work)
// 	if err != nil {
// 		utils.HandleError(c, http.StatusInternalServerError, "Failed to save work", err)
// 		return
// 	}
// 	// Update user's recent works
// 	user.CiWritten = append(user.CiWritten, user_work.ID)
// 	_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"token": token}, bson.M{"$set": bson.M{"recent_works": user.CiWritten}})
// 	if err != nil {
// 		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Work saved successfully",
// 		"work_id": user_work.ID.Hex(),
// 	})
// }

func PutIntoDraftsHandler(c *gin.Context, token string, draftObj models.ModernWork) {
	// 1. 先获取用户信息以获取用户ID
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)

	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// 2. 设置作者ID
	draftObj.Author = user.ID

	// 3. 将草稿存入drafts集合
	draftResult, err := config.MongoClient.Database("serenesong").Collection("drafts").InsertOne(c, draftObj)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoInsert, err)
		return
	}

	// 4. 获取插入的草稿ID
	draftID := draftResult.InsertedID.(primitive.ObjectID)

	// 5. 将草稿ID添加到用户的drafts数组
	result, err := config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c,
		bson.M{"token": token},
		bson.M{"$push": bson.M{"drafts": draftID}},
	)

	if err != nil {
		// 如果更新用户失败，需要删除已插入的草稿
		_, deleteErr := config.MongoClient.Database("serenesong").Collection("drafts").DeleteOne(c, bson.M{"_id": draftID})
		if deleteErr != nil {
			// 记录删除失败的错误，但不返回给用户
			log.Printf("Failed to delete draft after user update failed: %v", deleteErr)
		}
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}

	// 检查是否找到并更新了用户
	if result.MatchedCount == 0 {
		// 如果用户不存在，同样需要删除已插入的草稿
		_, deleteErr := config.MongoClient.Database("serenesong").Collection("drafts").DeleteOne(c, bson.M{"_id": draftID})
		if deleteErr != nil {
			log.Printf("Failed to delete draft after user not found: %v", deleteErr)
		}
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, nil)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Draft saved successfully", "draft_id": draftID.Hex()})
}

func DelDraftHandler(c *gin.Context, token string, draftID string) {
	// Find user by token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}
	draftIDObj, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid draft ID", err)
		return
	}
	// Find draft by ID
	for _, draft := range user.Drafts {
		if draft == draftIDObj {
			// Delete draft from user's drafts array
			_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"token": token}, bson.M{"$pull": bson.M{"drafts": draftIDObj}})
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
				return
			}
			// Delete draft from drafts collection
			_, err = config.MongoClient.Database("serenesong").Collection("drafts").DeleteOne(c, bson.M{"_id": draftIDObj})
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDelete, err)
				return
			}
			c.JSON(http.StatusOK, gin.H{
				"message": "Draft deleted successfully",
			})
			return
		}
	}
	// If draft not found, return error
	utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, nil)
}

func TurnToFormalHandler(c *gin.Context, token string, draftID string) {
	// Find user by token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// Find draft by ID
	draftIDObj, err := primitive.ObjectIDFromHex(draftID)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid draft ID", err)
		return
	}

	for _, draft := range user.Drafts {
		if draft == draftIDObj {
			// Find draft by ID
			var draftObj models.ModernWork
			err := config.MongoClient.Database("serenesong").Collection("drafts").FindOne(c, bson.M{"_id": draftIDObj}).Decode(&draftObj)
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
				return
			}

			// Save work to database
			draftObj.Author = user.ID
			draftObj.IsPublic = false
			insertResult, err := config.MongoClient.Database("serenesong").Collection("UserWorks").InsertOne(c, draftObj)
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoInsert, err)
				return
			}

			// Update user's recent works
			user.CiWritten = append(user.CiWritten, insertResult.InsertedID.(primitive.ObjectID))
			_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"token": token}, bson.M{"$set": bson.M{"ci_written": user.CiWritten}})
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
				return
			}

			// Delete draft from user's drafts array
			_, err = config.MongoClient.Database("serenesong").Collection("users").UpdateOne(c, bson.M{"token": token}, bson.M{"$pull": bson.M{"drafts": draftIDObj}})
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
				return
			}

			// Delete draft from drafts collection
			_, err = config.MongoClient.Database("serenesong").Collection("drafts").DeleteOne(c, bson.M{"_id": draftIDObj})
			if err != nil {
				utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDelete, err)
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message": "Work saved successfully",
				"work_id": insertResult.InsertedID.(primitive.ObjectID).Hex(),
			})
			return
		}
	}
	// If draft not found, return error
	utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, nil)
}

func ModifyModernWorkHandler(c *gin.Context, token string, workID string, work models.ModernWork, collectionName string) {
	// Find user by token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// Find work by ID
	workIDObj, err := primitive.ObjectIDFromHex(workID)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidObjID, err)
		return
	}

	// Check if user is the author of the work
	belongTo := false
	if collectionName == "UserWorks" {
		for _, workID := range user.CiWritten {
			if workID == workIDObj {
				belongTo = true
				break
			}
		}
	} else if collectionName == "drafts" {
		for _, workID := range user.Drafts {
			if workID == workIDObj {
				belongTo = true
				break
			}
		}
	}
	if !belongTo {
		utils.HandleError(c, http.StatusForbidden, utils.ErrMsgPermission, nil)
		return
	}

	// Update work
	filter := bson.M{"_id": workIDObj}
	update := bson.M{
		"$set": bson.M{
			"title":      work.Title,
			"content":    work.Content,
			"cipai":      work.Cipai,
			"prologue":   work.Xiaoxu,
			"is_public":  work.IsPublic,
			"tags":       work.Tags,
			"updated_at": time.Now(),
		},
	}

	result, err := config.MongoClient.Database("serenesong").Collection(collectionName).UpdateOne(c, filter, update)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoUpdate, err)
		return
	}
	if result.MatchedCount == 0 {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgMongoFind, nil)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Updated successfully"})
}

func GetMyWorksHandler(c *gin.Context, token string, collectionName string) {
	// Find user by token
	var user models.User
	err := config.MongoClient.Database("serenesong").Collection("users").FindOne(c, bson.M{"token": token}).Decode(&user)
	if err != nil {
		utils.HandleError(c, http.StatusNotFound, utils.ErrMsgUserNotFound, err)
		return
	}

	// Find drafts by user ID
	var drafts []models.ModernWork
	filter := bson.M{"author": user.ID}
	cursor, err := config.MongoClient.Database("serenesong").Collection(collectionName).Find(c, filter)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	defer cursor.Close(c)
	for cursor.Next(c) {
		var draft models.ModernWork
		err := cursor.Decode(&draft)
		if err != nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
			return
		}
		drafts = append(drafts, draft)
	}
	if err := cursor.Err(); err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoCursor, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		collectionName: drafts,
	})
}
