package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvValidator() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}

func EnvMongoUri() string {
	EnvValidator()
	return os.Getenv("MONGO_URI")
}

func SendGridKeyAndFrom() (string, string) {
	EnvValidator()
	return os.Getenv("SENDGRID_API_KEY"), os.Getenv("SENDGRID_SENDER")
}
