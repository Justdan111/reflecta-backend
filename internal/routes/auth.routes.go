package routes

import (
	"reflecta/internal/controllers"
	"reflecta/internal/middleware"

	"github.com/gofiber/fiber/v2"
)

func AuthRoutes(app *fiber.App) {
	api := app.Group("/api/auth")

	api.Post("/register", controllers.Register)
	api.Post("/login", controllers.Login)
	api.Get("/profile", middleware.AuthMiddleware, controllers.GetProfile)
}
