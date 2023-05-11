package tools

import (
	"log"

	"github.com/joho/godotenv"
)

func Load_env() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
