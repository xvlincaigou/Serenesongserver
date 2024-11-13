package controllers

import (
	"Serenesongserver/services"

	"github.com/gin-gonic/gin"
)

func CreateCollection(c *gin.Context) {
	collectionName := c.Query("collectionName")
	token := c.Query("token")
	if collectionName == "" || token == "" {
		c.JSON(400, gin.H{"error": "CollectionName and token are required"})
		return
	}
	services.CreateCollectionHandler(c, collectionName, token)
}

func DeleteCollection(c *gin.Context) {
	collectionID := c.Query("collectionID")
	token := c.Query("token")
	if collectionID == "" || token == "" {
		c.JSON(400, gin.H{"error": "CollectionID and token are required"})
		return
	}
	services.DeleteCollectionHandler(c, collectionID, token)
}

func AddToCollection(c *gin.Context) {
	collectionID := c.Query("collectionID")
	ciID := c.Query("ciID")
	token := c.Query("token")
	if collectionID == "" || ciID == "" || token == "" {
		c.JSON(400, gin.H{"error": "CollectionID, CiID and token are required"})
		return
	}
	services.AddToCollectionHandler(c, collectionID, ciID, token)
}

func RemoveFromCollection(c *gin.Context) {
	collectionID := c.Query("collectionID")
	ciID := c.Query("ciID")
	token := c.Query("token")
	if collectionID == "" || ciID == "" || token == "" {
		c.JSON(400, gin.H{"error": "CollectionID, CiID and token are required"})
		return
	}
	services.RemoveFromCollectionHandler(c, collectionID, ciID, token)
}

func GetAllColletions(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(400, gin.H{"error": "Token is required"})
		return
	}
	services.GetAllCollectionsHandler(c, token)
}

func GetAllColletionItems(c *gin.Context) {
	collectionID := c.Query("collectionID")
	token := c.Query("token")
	if collectionID == "" || token == "" {
		c.JSON(400, gin.H{"error": "CollectionID and token are required"})
		return
	}
	services.GetAllCollectionItemsHandler(c, collectionID, token)
}

func ModifyCollectionComment(c *gin.Context) {
	ciID := c.Query("ciID")
	comment := c.Query("comment")
	token := c.Query("token")
	if ciID == "" || comment == "" || token == "" {
		c.JSON(400, gin.H{"error": "CollectionID, Comment and token are required"})
		return
	}
	services.ModifyCollectionCommentHandler(c, ciID, comment, token)
}
