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

func EnvTokenSecrets() (access, refresh string) {
	EnvValidator()
	access = os.Getenv("ACCESS_TOKEN_SECRET")
	refresh = os.Getenv("REFRESH_TOKEN_SECRET")
	return
}

func EnvAWSCredentials() (accessKey, secretKey string) {
	EnvValidator()
	accessKey = os.Getenv("AWS_ACCESS_KEY_ID")
	secretKey = os.Getenv("AWS_SECRET_ACCESS_KEY")
	return
}
