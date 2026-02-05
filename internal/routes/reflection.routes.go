package routes

import (
	"github.com/gofiber/fiber/v2"
	
	"reflecta/internal/controllers"
	"reflecta/internal/middleware"
)

func ReflectionRoutes(app *fiber.App) {
	api := app.Group("/api/reflections", middleware.AuthMiddleware)

	api.Post("/", controllers.CreateReflection)
}
