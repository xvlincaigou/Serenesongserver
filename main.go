package main

import (
	"Serenesongserver/config"
	"Serenesongserver/controllers"
	// "Serenesongserver/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	// 只有在需要测试推荐任务的时候才使用这个功能
	// utils.SetupCronJobs()

	router := gin.Default()

	router.POST("/login", controllers.Login)

	router.GET("/getAllCollections", controllers.GetAllColletions)
	router.GET("/getAllColletionItems", controllers.GetAllColletionItems)
	router.POST("/createCollection", controllers.CreateCollection)
	router.POST("/deleteCollection", controllers.DeleteCollection)
	router.POST("/addToCollection", controllers.AddToCollection)
	router.POST("/removeFromCollection", controllers.RemoveFromCollection)
	router.POST("/modifyCollectionComment", controllers.ModifyCollectionComment)
	router.GET("/getCollectionItemCount", controllers.GetCollectionItemCount)

	router.GET("/recommendCi", controllers.RecommendCi)
	router.GET("/recommendPic", controllers.RecommendPic)

	// Searching related APIs
	router.GET("/search", controllers.SearchRouter)

	router.GET("/getCiById", controllers.GetCiById)

	router.Run("0.0.0.0:8080")
}
