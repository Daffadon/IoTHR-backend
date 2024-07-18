package server

import (
	"os"

	"github.com/joho/godotenv"
)

func Init() {
	if err := godotenv.Load(); err != nil {
		panic("Failed to load .env file!")
	}
	port := os.Getenv("APP_PORT")
	router := NewRouter()
	router.Run(":" + port)
}
