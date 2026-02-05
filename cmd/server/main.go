package main

import (
	"log"
	"reflecta/internal/config"
	"reflecta/internal/database"
	"reflecta/internal/routes"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// load environment variables
	config.LoadEnv()

	// connect to mongoDB
	database.ConnectMongo()

	app := fiber.New()

	// rate limiter
	app.Use(limiter.New(limiter.Config{
	Max:        100,
	Expiration: time.Minute,
}))

	// logger
app.Use(logger.New())

	// setup routes
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.ReflectionRoutes(app)

	log.Println("Server running on http://localhost:4000")
	app.Listen(":4000")
}
