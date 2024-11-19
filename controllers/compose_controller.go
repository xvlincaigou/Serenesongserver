package controllers

import (
	"Serenesongserver/services"
	"Serenesongserver/utils"

	"net/http"
	"strconv"
	// "fmt"

	"github.com/gin-gonic/gin"
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