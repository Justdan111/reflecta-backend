package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"reflecta/internal/database"
	"reflecta/internal/models"
)

func CreateReflection(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	var body struct {
		Mood int    `json:"mood"`
		Note string `json:"note"`
	}

	if err := c.BodyParser(&body); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	uid, _ := primitive.ObjectIDFromHex(userID)
	today := time.Now().Format("2006-01-02")

	collection := database.DB.Collection("reflections")

	// Prevent duplicate reflection per day
	count, _ := collection.CountDocuments(context.Background(), bson.M{
		"user_id": uid,
		"date":    today,
	})

	if count > 0 {
		return c.Status(400).JSON(fiber.Map{"error": "Reflection for today already exists"})
	}

	reflection := models.Reflection{
		UserID:    uid,
		Mood:      body.Mood,
		Note:      body.Note,
		Date:      time.Now(),
		CreatedAt: time.Now(),
	}

	_, err := collection.InsertOne(context.Background(), reflection)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to save reflection"})
	}

return c.JSON(fiber.Map{
		"message": "Reflection saved",
		"data":    reflection,
	})
}