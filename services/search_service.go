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

func AuthorConsult(c *gin.Context, keywords []string) {
	collection := config.MongoClient.Database("serenesong").Collection("Author")
	// Matching rule
	var layers []bson.M
	for _, keyword := range keywords {
		layers = append(layers, bson.M{
			"name": bson.M{
				"$regex": keyword,
				"$options": "i",
			},
		})
	}
	filter := bson.M{"$and": layers}
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

func CipaiConsult(c *gin.Context, keywords []string) {
	collection := config.MongoClient.Database("serenesong").Collection("CipaiList")
	// Matched cipai list
	match_list := make(map[string]interface{})
	// Serching all cipai and fetch names
	curson, err := collection.Find(c, bson.M{})
	if err!= nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoFind, err)
		return
	}
	// Searching each cipai
	for curson.Next(c) {
		var all_cipai bson.M
		if err := curson.Decode(&all_cipai); err!= nil {
			utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgMongoDecode, err)
			return
		}
		for cipai, content := range all_cipai {
			// Check if the cipai contains the keyword
			var contains bool = true
			for _, keyword := range keywords {
				if (!contains) { break }
				contains = strings.Contains(strings.ToLower(cipai), strings.ToLower(keyword))
			}
			if contains {
				match_list[cipai] = content
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{"matched_cipai": match_list})
}

func CiConsult(c *gin.Context, keywords []string) {
	collection := config.MongoClient.Database("serenesong").Collection("Ci")
	// Matching rule
	var layers []bson.M
	for _, keyword := range keywords {
		layers = append(layers, bson.M{
			"$or": []bson.M{
				{"author":  bson.M{"$regex": keyword, "$options": "i"}},
				{"title":   bson.M{"$regex": keyword, "$options": "i"}},
				{"content": bson.M{"$elemMatch": bson.M{"$regex": keyword, "$options": "i"}}},
				{"cipai":   bson.M{"$elemMatch": bson.M{"$regex": keyword, "$options": "i"}}},
			},
		})
	}
	filter := bson.M{"$and": layers}
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