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
	router.GET("/getCollectionItem", controllers.GetCollectionItem)

	router.GET("/recommendCi", controllers.RecommendCi)
	router.GET("/recommendPic", controllers.RecommendPic)

	// Searching related APIs
	router.GET("/search", controllers.SearchRouter)
	// Composing related APIs
	router.GET("/getRhymes", controllers.GetRhymes)
	router.GET("/getYunshu", controllers.GetYunshu)
	router.GET("/getFormat", controllers.GetFormat)
	// router.POST("/doneWork", controllers.FinishWork)
	router.POST("/putIntoDrafts", controllers.PutIntoDrafts)
	router.DELETE("/delDraft", controllers.DelDraft)
	router.POST("/turnToFormal", controllers.TurnToFormal)
	router.POST("/modifyDraft", controllers.ModifyDraft)
	router.POST("/modifyWork", controllers.ModifyWork)
	router.GET("/getMyWorks", controllers.GetMyWorks)
	router.GET("/getCiById", controllers.GetCiById)
	// User info related APIs
	router.GET("/getDynamics", controllers.GetDynamics)
	router.GET("/getCollections", controllers.GetCollections)
	router.GET("/getSubscribers", controllers.GetSubscribers)
	router.GET("/getSubscribedTo", controllers.GetSubscribedTo)
	router.GET("/getPublicWorks", controllers.GetPublicWorks)
	router.GET("/getWorks", controllers.GetWorks)
	router.GET("/getAvatar", controllers.GetAvatar)
	router.GET("/changePrivacy", controllers.ChangePrivacy)

	router.Run("0.0.0.0:8080")
}
