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
