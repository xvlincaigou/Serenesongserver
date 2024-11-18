package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCollection(c *gin.Context) {
	collectionName := c.Query("collectionName")
	token := c.Query("token")
	if collectionName == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.CreateCollectionHandler(c, collectionName, token)
}

func DeleteCollection(c *gin.Context) {
	collectionID := c.Query("collectionID")
	token := c.Query("token")
	if collectionID == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.DeleteCollectionHandler(c, collectionID, token)
}

func AddToCollection(c *gin.Context) {
	collectionID := c.Query("collectionID")
	ciID := c.Query("ciID")
	token := c.Query("token")
	if collectionID == "" || ciID == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.AddToCollectionHandler(c, collectionID, ciID, token)
}

func RemoveFromCollection(c *gin.Context) {
	collectionID := c.Query("collectionID")
	ciID := c.Query("ciID")
	token := c.Query("token")
	if collectionID == "" || ciID == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.RemoveFromCollectionHandler(c, collectionID, ciID, token)
}

func GetAllColletions(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.GetAllCollectionsHandler(c, token)
}

func GetAllColletionItems(c *gin.Context) {
	collectionID := c.Query("collectionID")
	token := c.Query("token")
	if collectionID == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.GetAllCollectionItemsHandler(c, collectionID, token)
}

func ModifyCollectionComment(c *gin.Context) {
	ciID := c.Query("ciID")
	comment := c.Query("comment")
	token := c.Query("token")
	if ciID == "" || comment == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ModifyCollectionCommentHandler(c, ciID, comment, token)
}

func GetCollectionItemCount(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.GetCollectionItemCountHandler(c, token)
}
