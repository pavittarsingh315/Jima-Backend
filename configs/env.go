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

func EnvSendGridKeyAndFrom() (key, sender string) {
	EnvValidator()
	key = os.Getenv("SENDGRID_API_KEY")
	sender = os.Getenv("SENDGRID_SENDER")
	return
}

func EnvTwilioIDKeyFrom() (id, token, from string) {
	EnvValidator()
	id = os.Getenv("TWILIO_ACCOUNT_SID")
	token = os.Getenv("TWILIO_AUTH_TOKEN")
	from = os.Getenv("TWILIO_FROM_NUMBER")
	return
}
