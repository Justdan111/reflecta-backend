package controllers

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"reflecta/internal/database"
	"reflecta/internal/models"
	"reflecta/internal/utils"
)

func CreateReflection(c *fiber.Ctx) error {
	utils.ReflectionMutex.Lock()
	defer utils.ReflectionMutex.Unlock()

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

	if err := utils.ValidateMood(body.Mood); err != nil {
	return c.Status(400).JSON(fiber.Map{
		"error": err.Error(),
	})
}

return c.JSON(fiber.Map{
		"message": "Reflection saved",
		"data":    reflection,
	})
}

func WeeklySummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	uid, _ := primitive.ObjectIDFromHex(userID)

	start := time.Now().AddDate(0, 0, -7)

	collection := database.DB.Collection("reflections")

	cursor, err := collection.Find(context.Background(), bson.M{
		"user_id": uid,
		"created_at": bson.M{
			"$gte": start,
		},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Failed to fetch data"})
	}

	var reflections []models.Reflection
	cursor.All(context.Background(), &reflections)

	totalMood := 0
	for _, r := range reflections {
		totalMood += r.Mood
	}

	avgMood := 0
	if len(reflections) > 0 {
		avgMood = totalMood / len(reflections)
	}
	
	return c.JSON(fiber.Map{
		"count":    len(reflections),
		"avgMood": avgMood,
		"entries": reflections,
	})
}

func PersonalInsights(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	uid, _ := primitive.ObjectIDFromHex(userID)
	
	collection := database.DB.Collection("reflections")

	cursor, _ := collection.Find(context.Background(), bson.M{"user_id": uid})

	var reflections []models.Reflection
	cursor.All(context.Background(), &reflections)

	lowMoodDays := 0
	for _, r := range reflections {
		if r.Mood <= 2 {
			lowMoodDays++
		}
	}

	insight := "You are emotionally balanced recently 🌱"
	if lowMoodDays >= 3 {
		insight = "You’ve had several low days. Consider rest or reflection 💭"
	}

	return c.JSON(fiber.Map{
		"insight": insight,
	})
}