package controllers

import (
	"Serenesongserver/services"
	"github.com/gin-gonic/gin"
)

func RecommendCi(c *gin.Context) {
	services.RecommendCiHandler(c)
}

func RecommendPic(c *gin.Context) {
	services.RecommendPicHandler(c)
}
