package main

import (
	"Serenesongserver/config"
	"Serenesongserver/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	router := gin.Default()

	router.POST("/login", controllers.Login)

	router.GET("/getAllCollections", controllers.GetAllColletions)
	router.GET("/getAllColletionItems", controllers.GetAllColletionItems)
	router.POST("/createCollection", controllers.CreateCollection)
	router.POST("/deleteCollection", controllers.DeleteCollection)
	router.POST("/addToCollection", controllers.AddToCollection)
	router.POST("/removeFromCollection", controllers.RemoveFromCollection)
	router.POST("/modifyCollectionComment", controllers.ModifyCollectionComment)

	router.Run(":8080")
}
