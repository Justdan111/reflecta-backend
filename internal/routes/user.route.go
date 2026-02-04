package routes

import (
	"github.com/gofiber/fiber/v2"
	"reflecta/internal/middleware"
)

func UserRoutes(app *fiber.App) {
	api := app.Group("/api/user")

	api.Get("/me", middleware.AuthMiddleware, func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{"user_id": userID})
	})
}
