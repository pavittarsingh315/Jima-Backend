package main

import (
	"NeraJima/configs"
	"NeraJima/routes"
	"context"
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/helmet/v2"
)

func main() {
	configs.EnvValidator()

	app := fiber.New()

	// Middleware
	app.Use(helmet.New())
	app.Use(cors.New())
	app.Use(logger.New())

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := configs.ConnectDatabase()
	defer client.Disconnect(ctx)
	routes.SetupRouter(app)

	port := os.Getenv("PORT")
	err := app.Listen(":" + port)

	if err != nil {
		log.Fatal("ERROR: app failed to start")
		panic(err)
	}
}
