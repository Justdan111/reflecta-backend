package main

import (
	"log"
	

	"github.com/gofiber/fiber/v2"
	"refecta/config"
	"refecta/database"
	"refecta/routes"
	
)

func main() {
	// load environment variables
	config.LoadEnv()

	// connect to mongoDB
	database.ConnectMongo()

	app := fiber.New()

	// setup routes
	routes.AuthRoutes(app)

	log.Println("Server running on http://localhost:4000")
	app.Listen(":4000")
}