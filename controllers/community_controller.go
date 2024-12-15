package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

// 发布动态
func PublishDynamic(c *gin.Context) {
	var jsonData bson.M
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}

	token, tokenOk := jsonData["token"].(string)
	_Type, _TypeOk := jsonData["Type"].(float64)
	_id, _idOk := jsonData["_id"].(string)
	_id2, _id2Ok := jsonData["_id2"].(string)

	if !tokenOk || !_TypeOk || !_idOk || !_id2Ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	Type := int(_Type)
	services.PublishDynamicHandler(c, token, Type, _id, _id2)
}

func WithdrawPost(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}
	token, token_ok := json_data["token"].(string)
	post_id, post_ok := json_data["post_id"].(string)
	if !token_ok || !post_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.DeletePost(c, token, post_id)
}

func GetRandomPosts(c *gin.Context) {
	token := c.Query("token")
	v_str := c.Query("value")
	if token == "" || v_str == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	// Check if value is not an integer
	value, err := strconv.Atoi(v_str)
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgNotAnInteger, nil)
		return
	}
	if value <= 0 {
		utils.HandleError(c, http.StatusBadRequest, "Must be greater than 0", nil)
		return
	}
	services.ReturnRandomPosts(c, token, value)
}

func GetFollowingPosts(c *gin.Context) {
	user_id := c.Query("user_id")
	token := c.Query("token")
	if user_id == "" || token == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.ReturnFollowingPosts(c, user_id, token)
}

func CommentPost(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}
	token, token_ok := json_data["token"].(string)
	content, content_ok := json_data["content"].(string)
	post_id, post_ok := json_data["post_id"].(string)
	if !token_ok || !content_ok || !post_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.SendComment(c, token, content, post_id)
}

func WithdrawComment(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}
	token, token_ok := json_data["token"].(string)
	comment_id, comment_ok := json_data["comment_id"].(string)
	if !token_ok || !comment_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.DeleteComment(c, token, comment_id)
}

func LikePost(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}
	token, token_ok := json_data["token"].(string)
	post_id, post_ok := json_data["post_id"].(string)
	if !token_ok || !post_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.SendLike(c, token, post_id)
}

func WithdrawLike(c *gin.Context) {
	var json_data bson.M
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, err)
		return
	}
	token, token_ok := json_data["token"].(string)
	post_id, post_ok := json_data["post_id"].(string)
	if !token_ok || !post_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.DeleteLike(c, token, post_id)
}
