package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"net/http"
	"strconv"
	"encoding/json"
	// "fmt"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	// "go.mongodb.org/mongo-driver/mongo/options"
)

func GetRhymes(c *gin.Context) {
	services.ReturnRhymes(c)
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
	// Get the new work & the token from the query string
	new_work := c.Query("new_work")
	token := c.Query("token")
	// Extract the work data to JSON
	var work_data bson.M
	err := json.Unmarshal([]byte(new_work), &work_data)
	if err!= nil {
		utils.HandleError(c, http.StatusBadRequest, utils.ErrMsgInvalidJSON, nil)
		return
	}
	services.SaveWork(c, work_data, token)
}