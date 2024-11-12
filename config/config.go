package config

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	AppID       string
	Secret      string
	MongoClient *mongo.Client
)

func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	AppID = os.Getenv("APPID")
	Secret = os.Getenv("SECRET")
	MongoDBURI := os.Getenv("MONGODBURI")

	clientOptions := options.Client().ApplyURI(MongoDBURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	MongoClient = client
}
