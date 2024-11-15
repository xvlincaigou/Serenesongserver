package services

import (
	"Serenesongserver/utils"
	"github.com/gin-gonic/gin"
	"net/http"
)

func RecommendCiHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"recommend": utils.RecommendedCi})
}

func RecommendPicHandler(c *gin.Context) {
	if utils.RecommendedPicPath != "" {
		c.File(utils.RecommendedPicPath)
	}
}
