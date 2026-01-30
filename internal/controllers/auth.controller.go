package controllers

import (
	"context"
	"reflecta/internal/database"
	"reflecta/internal/models"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var body models.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	body.Password = string(hash)

	collection := database.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := collection.InsertOne(ctx, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "User creation failed"})
	}

	return c.JSON(fiber.Map{"message": "User registered"})
}

func Login(c *fiber.Ctx) error {
	var body models.User
	c.BodyParser(&body)

	collection := database.DB.Collection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&user)

	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	return c.JSON(fiber.Map{"message": "Login successful"})
}
