package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(200, "Hello World!")
	})
	router.GET("/bye", func(c *gin.Context) {
		c.String(200, "Goodbye World!")
	})
	router.Run(":8080")
}