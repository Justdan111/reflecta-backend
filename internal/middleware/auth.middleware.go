package middleware

import (
	"strings"

	"reflecta/internal/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *fiber.Ctx) error {
   authHeader := c.Get("Authorization")

   if authHeader == "" {
	  return c.Status(401).JSON(fiber.Map{"error": "Missing token"})
   }

   tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

   secret := config.GetEnv("JWT_SECRET")

   token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
	  return []byte(secret), nil
	     })

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid token"})
	}
	
	claims := token.Claims.(jwt.MapClaims)
	c.Locals("user_id", claims["user_id"])
	
	return c.Next()

}