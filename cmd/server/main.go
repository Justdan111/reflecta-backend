package main

import (
	"log"
	"reflecta/internal/config"
	"reflecta/internal/database"
	"reflecta/internal/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	// load environment variables
	config.LoadEnv()

	// connect to mongoDB
	database.ConnectMongo()

	app := fiber.New()

	// setup routes
	routes.AuthRoutes(app)
	routes.UserRoutes(app)
	routes.ReflectionRoutes(app)

	log.Println("Server running on http://localhost:4000")
	app.Listen(":4000")
}
