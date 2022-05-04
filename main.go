package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
	"github.com/joho/godotenv"
)

func main() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	app := fiber.New()

	// Middleware
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello World 👋! Its me, PSJ")
	})

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("ERROR: app failed to start")
		panic(err)
	}
}
