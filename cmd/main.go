package main

import (
	"app/cmd/auth/controller"
	"app/internal/database"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using default environment variables")
	}
	database_config := database.Parameters{
		Host:         os.Getenv("DB_HOST"),
		Port:         os.Getenv("DB_PORT"),
		User:         os.Getenv("DB_USER"),
		Password:     os.Getenv("DB_PASSWORD"),
		DatabaseName: os.Getenv("DB_NAME"),
	}
	database.Init(database_config)

	router := gin.Default()

	router.POST("/auth/get", controller.GetTokens)
	router.POST("/auth/refresh", controller.RefreshTokens)
	router.Run(":8000")
}
