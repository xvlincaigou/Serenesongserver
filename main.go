package main

import (
	"Serenesongserver/routers"
	"os"
	"log"
	"context"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// 从 miniapp.env 文件中加载环境变量
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	appid := os.Getenv("APPID")
	secret := os.Getenv("SECRET")

	// 初始化mongodb
	mongodburi := os.Getenv("MONGODBURI")
	clientOptions := options.Client().ApplyURI(mongodburi)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	router := gin.Default()
	router.GET("/login", func(c *gin.Context) {
		routers.LoginHandler(c, appid, secret, client)
	})
	router.Run(":8080")
}
