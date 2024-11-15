package main

import (
	"Serenesongserver/config"
	"Serenesongserver/controllers"
	"Serenesongserver/utils"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
)

func main() {
	config.LoadEnv()

	c := cron.New()
	c.Start()
	c.AddFunc("@every 24h", utils.GenerateAndDownloadImageWrapper)
	defer c.Stop()

	router := gin.Default()

	router.POST("/login", controllers.Login)

	router.GET("/getAllCollections", controllers.GetAllColletions)
	router.GET("/getAllColletionItems", controllers.GetAllColletionItems)
	router.POST("/createCollection", controllers.CreateCollection)
	router.POST("/deleteCollection", controllers.DeleteCollection)
	router.POST("/addToCollection", controllers.AddToCollection)
	router.POST("/removeFromCollection", controllers.RemoveFromCollection)
	router.POST("/modifyCollectionComment", controllers.ModifyCollectionComment)

	router.GET("/recommendCi", controllers.RecommendCi)
	router.GET("/recommendPic", controllers.RecommendPic)

	router.Run("0.0.0.0:8080")
}
