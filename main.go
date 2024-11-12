package main

import (
	"Serenesongserver/config"
	"Serenesongserver/controllers"

	"github.com/gin-gonic/gin"
)

func main() {
	config.LoadEnv()

	router := gin.Default()

	router.GET("/login", controllers.Login)

	router.POST("/createCollection", controllers.CreateCollection)
	router.POST("/deleteCollection", controllers.DeleteCollection)
	router.POST("/addToCollection", controllers.AddToCollection)
	router.POST("/removeFromCollection", controllers.RemoveFromCollection)

	router.Run(":8080")
}
