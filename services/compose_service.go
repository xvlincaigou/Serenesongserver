package services

import (
	// "strings"
	// "fmt"

	"Serenesongserver/config"
	// "Serenesongserver/models"
	"Serenesongserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "golang.org/x/text/message/pipeline"
	// "go.mongodb.org/mongo-driver/bson/primitive"
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

	c.JSON(http.StatusOK, gin.H{"rhymes": rhymes, "pingze": pingze})
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