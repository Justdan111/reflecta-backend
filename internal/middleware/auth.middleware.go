package middleware

import (
	"context"
	"strings"
	"time"

	"reflecta/internal/config"
	"reflecta/internal/database"
	"reflecta/internal/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
)

func AuthMiddleware(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return c.Status(401).JSON(fiber.Map{"message": "Missing token"})
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	secret := config.GetEnv("JWT_SECRET")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "Session expired. Please log in again."})
	}

	// Reject tokens that were revoked by logout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := database.DB.Collection("revoked_tokens").CountDocuments(ctx, bson.M{"token_hash": utils.HashToken(tokenString)})
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Could not verify session"})
	}
	if count > 0 {
		return c.Status(401).JSON(fiber.Map{"message": "Session expired. Please log in again."})
	}

	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", claims["user_id"])
	c.Locals("token", tokenString)
	c.Locals("claims", claims)

	return c.Next()
}
