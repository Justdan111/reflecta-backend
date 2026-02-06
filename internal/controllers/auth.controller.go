package controllers

import (
	"context"
	"reflecta/internal/database"
	"reflecta/internal/models"
	"reflecta/internal/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *fiber.Ctx) error {
	var body models.User
	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid input"})
	}

	collection := database.DB.Collection("users")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Check if email already exists
	var existingUser models.User
	err := collection.FindOne(ctx, bson.M{"email": body.Email}).Decode(&existingUser)
	if err == nil {
		return c.Status(400).JSON(fiber.Map{"message": "Email already exists"})
	}

	// Hash password
	hash, _ := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	body.Password = string(hash)
	body.ID = primitive.NewObjectID()
	body.CreatedAt = time.Now().Unix()

	_, err = collection.InsertOne(ctx, body)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "User creation failed"})
	}

	// Generate JWT token
	token, _ := utils.GenerateToken(body.ID.Hex())

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    body.ID.Hex(),
			"name":  body.Name,
			"email": body.Email,
		},
	})
}

func Login(c *fiber.Ctx) error {
	var body models.User
	c.BodyParser(&body)

	collection := database.DB.Collection("users")

	var user models.User
	err := collection.FindOne(context.Background(), bson.M{"email": body.Email}).Decode(&user)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid credentials"})
	}

	// Generate JWT token
	token, _ := utils.GenerateToken(user.ID.Hex())

	return c.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID.Hex(),
			"name":  user.Name,
			"email": user.Email,
		},
	})
}

func GetProfile(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	objectID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"message": "Invalid user ID"})
	}

	collection := database.DB.Collection("users")

	var user models.User
	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&user)

	if err != nil {
		return c.Status(404).JSON(fiber.Map{"message": "User not found"})
	}

	return c.JSON(fiber.Map{
		"id":    user.ID.Hex(),
		"name":  user.Name,
		"email": user.Email,
	})
}
