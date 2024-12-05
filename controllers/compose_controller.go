package controllers

import (
	"Serenesongserver/models"
	"Serenesongserver/services"
	"Serenesongserver/utils"

	// "encoding/json"
	"net/http"
	"strconv"
	// "encoding/json"
	// "fmt"

	"github.com/gin-gonic/gin"
	// "go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func GetRhymes(c *gin.Context) {
	services.ReturnRhymes(c)
}

func GetYunshu(c *gin.Context) {
	services.ReturnYunshu(c)
}

func GetFormat(c *gin.Context) {
	// Turn the format string into an integer
	cipai_name := c.Query("cipai_name")
	format_num, err := strconv.Atoi(c.Query("format_num"))
	if err != nil {
		utils.HandleError(c, http.StatusBadRequest, "Invalid format number", nil)
		return
	}
	services.ReturnFormat(c, cipai_name, format_num)
}

func FinishWork(c *gin.Context) {
	// // Get the new work & the token from the query string
	// new_work := c.Query("new_work")
	// token := c.Query("token")
	// Get the data from the request body.
	var json_data map[string]interface{}
	if err := c.ShouldBindJSON(&json_data); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	// Extract the token and the draft data from the JSON.
	token, token_ok := json_data["token"].(string)
	work, work_ok := json_data["new_work"].(map[string]interface{})
	if !token_ok || !work_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	// // Extract the work data to JSON
	// var work_data bson.M
	// err := json.Unmarshal([]byte(new_work), &work_data)
	// if err!= nil {
	// 	utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
	// 	return
	// }
	services.SaveWork(c, work, token)
}

func PutIntoDrafts(c *gin.Context) {
	// Get the data from the request body.
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Extract the token and the draft data from the JSON.
	token, token_ok := jsonData["token"].(string)
	draft, draft_ok := jsonData["draft"].(map[string]interface{})
	if !token_ok || !draft_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Create a new draft object from the data.
	draftObj, err := models.NewModernWork(draft)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	}

	services.PutIntoDraftsHandler(c, token, draftObj)
}

func DelDraft(c *gin.Context) {
	// Get the data from the request body.
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Extract the token and the draftid from the JSON.
	token, token_ok := jsonData["token"].(string)
	draftID, draftID_ok := jsonData["draftID"].(string)
	if !token_ok || !draftID_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	services.DelDraftHandler(c, token, draftID)
}

func TurnToFormal(c *gin.Context) {
	// Get the data from the request body.
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Extract the token and the draftid from the JSON.
	token, token_ok := jsonData["token"].(string)
	draftID, draftID_ok := jsonData["draftID"].(string)
	if !token_ok || !draftID_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	services.TurnToFormalHandler(c, token, draftID)
}

func ModifyDraft(c *gin.Context) {
	// Get the data from the request body.
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Extract the token and the draftid and draft data from the JSON.
	token, token_ok := jsonData["token"].(string)
	draftID, draftID_ok := jsonData["draftID"].(string)
	draft, draft_ok := jsonData["draft"].(map[string]interface{})
	if !token_ok || !draftID_ok || !draft_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Create a new draft object from the data.
	draftObj, err := models.NewModernWork(draft)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	}

	services.ModifyModernWorkHandler(c, token, draftID, draftObj, "drafts") // 这里的collectionName不是收藏夹名字，是mongodb的collection名字
}

func ModifyWork(c *gin.Context) {
	// Get the data from the request body.
	var jsonData map[string]interface{}
	if err := c.ShouldBindJSON(&jsonData); err != nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Extract the token and the workid and work data from the JSON.
	token, token_ok := jsonData["token"].(string)
	workID, workID_ok := jsonData["workID"].(string)
	work, work_ok := jsonData["work"].(map[string]interface{})
	if !token_ok || !workID_ok || !work_ok {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}

	// Create a new work object from the data.
	workObj, err := models.NewModernWork(work)
	if err != nil {
		utils.HandleError(c, http.StatusInternalServerError, utils.ErrMsgInternalError, err)
		return
	}

	services.ModifyModernWorkHandler(c, token, workID, workObj, "UserWorks") // 这里的collectionName不是收藏夹名字，是mongodb的collection名字
}

func GetMyWorks(c *gin.Context) {
	// Get the token from the query string.
	token := c.Query("token")
	kind := c.Query("kind")
	if token == "" || kind == "" {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidParams, nil)
		return
	}
	services.GetMyWorksHandler(c, token, kind)
}
