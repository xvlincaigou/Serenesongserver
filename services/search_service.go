package services

import (
	"strings"

	"Serenesongserver/config"
	"Serenesongserver/models"
	"Serenesongserver/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "golang.org/x/text/message/pipeline"
	// "go.mongodb.org/mongo-driver/bson/primitive"
)

func AuthorConsult(c *gin.Context, keyword string) {
	collection := config.MongoClient.Database("serenesong").Collection("Author")
	// Matching rule
	filter := bson.M{
		"name": bson.M{
			"$regex": keyword,
			"$options": "i",
		},
	}
	// Applying the filter and sorting
	cursor, err := collection.Find(c, filter)
	if err!= nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	var authors []models.Author
	if err := cursor.All(c, &authors); err!= nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"matched_authors": authors})
}

func CipaiConsult(c *gin.Context, keyword string) {
	collection := config.MongoClient.Database("serenesong").Collection("CipaiList")
	// Matched cipai list
	match_list := make(map[string]interface{})
	// Serching all cipai and fetch names
	curson, err := collection.Find(c, bson.M{})
	if err!= nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	for curson.Next(c) {
		var all_cipai bson.M
		if err := curson.Decode(&all_cipai); err!= nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
			return
		}
		for cipai, content := range all_cipai {
			if strings.Contains(strings.ToLower(cipai), strings.ToLower(keyword)) {
				match_list[cipai] = content
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"matched_cipai": match_list})
}

func CiConsult(c *gin.Context, keyword string) {
	collection := config.MongoClient.Database("serenesong").Collection("Ci")
	// Matching rule
	filter := bson.M{
		"$or": []bson.M{
			{"author":  bson.M{"$regex": keyword, "$options": "i"}},
			{"title":   bson.M{"$regex": keyword, "$options": "i"}},
			{"content": bson.M{"$elemMatch": bson.M{"$regex": keyword, "$options": "i"}}},
			{"cipai":   bson.M{"$elemMatch": bson.M{"$regex": keyword, "$options": "i"}}},
		},
	}
	// Applying the filter and sorting
	cursor, err := collection.Find(c, filter)
	if err!= nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	var ci []models.Ci
	if err := cursor.All(c, &ci); err!= nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"matched_ci": ci})
}