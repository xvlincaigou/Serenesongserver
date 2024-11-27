package controllers

import (
	// "Serenesongserver/models"
	"Serenesongserver/services"
	"Serenesongserver/utils"
	// "go/token"

	// "encoding/json"
	"net/http"
	// "strconv"
	// "encoding/json"
	// "fmt"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func GetDynamics(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user_id == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnDynamics(c, user_id, token)
}

func GetCollections(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user_id == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnCollections(c, user_id, token)
}

func GetSubscribers(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnSubscribers(c, token)
}

func GetSubscribedTo(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnSubscribedTo(c, token)
}