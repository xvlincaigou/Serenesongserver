package controllers

import (
	// "Serenesongserver/models"
	"Serenesongserver/services"
	"Serenesongserver/utils"
	// "os/user"

	// "go/token"

	// "encoding/json"
	"net/http"
	"strconv"

	// "encoding/json"
	// "fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

func GetPublicWorks(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user_id == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnPublicWorks(c, user_id, token)
}

func GetWorks(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user_id == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnWorks(c, user_id, token)
}

func GetUserInfo(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user_id == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnUserInfo(c, user_id, token)
}

func ChangePrivacy(c *gin.Context) {
	// Get the data from the request body.
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	work_id, work_ok    := json_data["work_id"].(string)
	token, token_ok 	:= json_data["token"].(string)
	privacy, privacy_ok := json_data["privacy"].(string)
	if !work_ok || !token_ok || !privacy_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	// Parse the 'privacy' parameter to a boolean
    is_public, err := strconv.ParseBool(privacy)
    if err != nil {
        utils.HandleError(c, http.StatusBadRequest, "Invalid privacy value. Must be a boolean.", nil)
        return
    }
	services.ChangePrivacy(c, work_id, token, is_public)
}

func SaveUserInfo(c *gin.Context) {
	// Get the data from the request body.
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	// Extract the token and the draft data from the JSON.
	token, token_ok   := json_data["token"].(string)
	avatar, avatar_ok := json_data["avatar"].(string)
	name, name_ok     := json_data["name"].(string)
	signature, sig_ok := json_data["signature"].(string)
	if !token_ok || !avatar_ok || !name_ok || !sig_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	services.SaveNameAvatar(c, token, name, avatar, signature)
}