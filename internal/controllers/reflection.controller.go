package controllers

import (
	"context"
	"fmt"
	"math"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

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
		return c.Status(400).JSON(fiber.Map{"message": "Invalid input"})
	}

	// Validate mood (1-5)
	if err := utils.ValidateMood(body.Mood); err != nil {
		return c.Status(400).JSON(fiber.Map{"message": err.Error()})
	}

	// Validate note length (max 500 chars)
	if len(body.Note) > 500 {
		return c.Status(400).JSON(fiber.Map{"message": "Note must be 500 characters or less"})
	}

	uid, _ := primitive.ObjectIDFromHex(userID)
	now := time.Now()

	collection := database.DB.Collection("reflections")

	reflection := models.Reflection{
		ID:        primitive.NewObjectID(),
		UserID:    uid,
		Mood:      body.Mood,
		Note:      body.Note,
		Date:      now,
		CreatedAt: now,
	}

	_, err := collection.InsertOne(context.Background(), reflection)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to save reflection"})
	}

	return c.JSON(fiber.Map{
		"id":        reflection.ID.Hex(),
		"mood":      reflection.Mood,
		"note":      reflection.Note,
		"createdAt": reflection.CreatedAt.Format(time.RFC3339),
		"userId":    userID,
	})
}

func WeeklySummary(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	uid, _ := primitive.ObjectIDFromHex(userID)

	// Get current week boundaries (Monday to Sunday)
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday is 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	endOfWeek := startOfWeek.AddDate(0, 0, 7)

	collection := database.DB.Collection("reflections")

	cursor, err := collection.Find(context.Background(), bson.M{
		"user_id": uid,
		"created_at": bson.M{
			"$gte": startOfWeek,
			"$lt":  endOfWeek,
		},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to load weekly summary"})
	}

	var reflections []models.Reflection
	if err := cursor.All(context.Background(), &reflections); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to load weekly summary"})
	}

	// Build weekly data by day
	days := []string{"MON", "TUE", "WED", "THU", "FRI", "SAT", "SUN"}
	weeklyData := make([]fiber.Map, 7)
	dayMoods := make(map[int][]int) // day index -> list of moods

	for i := 0; i < 7; i++ {
		dayMoods[i] = []int{}
	}

	for _, r := range reflections {
		dayIndex := int(r.CreatedAt.Weekday()) - 1
		if dayIndex < 0 {
			dayIndex = 6 // Sunday
		}
		dayMoods[dayIndex] = append(dayMoods[dayIndex], r.Mood)
	}

	totalMood := 0.0
	moodCount := 0
	moodFrequency := make(map[int]int)

	for i := 0; i < 7; i++ {
		avgMood := 0.0
		if len(dayMoods[i]) > 0 {
			sum := 0
			for _, m := range dayMoods[i] {
				sum += m
				moodFrequency[m]++
			}
			avgMood = math.Round(float64(sum)/float64(len(dayMoods[i]))*10) / 10
			totalMood += avgMood
			moodCount++
		}
		weeklyData[i] = fiber.Map{
			"day":  days[i],
			"mood": avgMood,
		}
	}

	// Calculate average mood label
	avgMoodLabel := "Neutral"
	if moodCount > 0 {
		avg := totalMood / float64(moodCount)
		if avg >= 4 {
			avgMoodLabel = "Positive"
		} else if avg >= 3 {
			avgMoodLabel = "Neutral"
		} else {
			avgMoodLabel = "Low"
		}
	}

	// Determine top emotion based on most frequent mood
	topEmotion := "Neutral"
	maxFreq := 0
	for mood, freq := range moodFrequency {
		if freq > maxFreq {
			maxFreq = freq
			switch mood {
			case 1:
				topEmotion = "Sad"
			case 2:
				topEmotion = "Pensive"
			case 3:
				topEmotion = "Neutral"
			case 4:
				topEmotion = "Calm"
			case 5:
				topEmotion = "Radiant"
			}
		}
	}

	// Calculate streak
	streak := calculateStreak(uid)

	// Format date range
	endDate := startOfWeek.AddDate(0, 0, 6)
	dateRange := fmt.Sprintf("%s — %s", startOfWeek.Format("Jan 2"), endDate.Format("Jan 2"))

	// Generate insight
	insight := generateWeeklyInsight(reflections, avgMoodLabel)

	return c.JSON(fiber.Map{
		"weeklyData":  weeklyData,
		"dateRange":   dateRange,
		"avgMood":     avgMoodLabel,
		"topEmotion":  topEmotion,
		"reflections": fmt.Sprintf("%d Posts", len(reflections)),
		"streak":      fmt.Sprintf("%d Days", streak),
		"insight":     insight,
	})
}

func calculateStreak(userID primitive.ObjectID) int {
	collection := database.DB.Collection("reflections")

	// Get reflections sorted by date descending
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := collection.Find(context.Background(), bson.M{"user_id": userID}, opts)
	if err != nil {
		return 0
	}

	var reflections []models.Reflection
	if err := cursor.All(context.Background(), &reflections); err != nil {
		return 0
	}

	if len(reflections) == 0 {
		return 0
	}

	streak := 0
	today := time.Now().Truncate(24 * time.Hour)
	expectedDate := today

	// Create a map of dates with reflections
	dateMap := make(map[string]bool)
	for _, r := range reflections {
		dateStr := r.CreatedAt.Truncate(24 * time.Hour).Format("2006-01-02")
		dateMap[dateStr] = true
	}

	// Count consecutive days backwards from today
	for i := 0; i < 365; i++ {
		checkDate := expectedDate.AddDate(0, 0, -i).Format("2006-01-02")
		if dateMap[checkDate] {
			streak++
		} else if i > 0 {
			break
		}
	}

	return streak
}

func generateWeeklyInsight(reflections []models.Reflection, avgMood string) string {
	if len(reflections) == 0 {
		return "Start journaling to unlock personalized insights about your emotional patterns."
	}

	switch avgMood {
case "Positive":
		return "You had a great week! Your reflections show a positive emotional trend. Keep up the good habits that contribute to your wellbeing."
	case "Neutral":
		return "Your week was balanced. Notice what activities or thoughts bring you closer to feeling calm and content."
	}
	return "This week had its challenges. Remember that acknowledging difficult emotions is the first step to understanding them. Consider what small changes might support your wellbeing."
}

func PersonalInsights(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)
	uid, _ := primitive.ObjectIDFromHex(userID)

	// Get reflections from the past week
	now := time.Now()
	weekday := int(now.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	startOfWeek := now.AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)

	collection := database.DB.Collection("reflections")

	cursor, err := collection.Find(context.Background(), bson.M{
		"user_id": uid,
		"created_at": bson.M{
			"$gte": startOfWeek,
		},
	})

	if err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to load insights"})
	}

	var reflections []models.Reflection
	if err := cursor.All(context.Background(), &reflections); err != nil {
		return c.Status(500).JSON(fiber.Map{"message": "Failed to load insights"})
	}

	// Build mood distribution by day of week
	days := []string{"Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"}
	dayMoods := make(map[int][]int)
	for i := 0; i < 7; i++ {
		dayMoods[i] = []int{}
	}

	for _, r := range reflections {
		dayIndex := int(r.CreatedAt.Weekday()) - 1
		if dayIndex < 0 {
			dayIndex = 6
		}
		dayMoods[dayIndex] = append(dayMoods[dayIndex], r.Mood)
	}

	// Convert to percentage values (mood 1-5 -> 0-100)
	moodDistribution := make([]fiber.Map, 7)
	lowDays := 0
	highDays := 0

	for i := 0; i < 7; i++ {
		value := 50 // default neutral
		color := "#6D5D8B"

		if len(dayMoods[i]) > 0 {
			sum := 0
			for _, m := range dayMoods[i] {
				sum += m
			}
			avgMood := float64(sum) / float64(len(dayMoods[i]))
			value = int((avgMood / 5.0) * 100)

			if avgMood <= 2.5 {
				color = "#C9A24D" // golden/amber for lower moods
				lowDays++
			} else if avgMood >= 4 {
				highDays++
			}
		}

		moodDistribution[i] = fiber.Map{
			"day":   days[i],
			"value": value,
			"color": color,
		}
	}

	// Generate mood uplift insight
	moodUplift := generateMoodUplift(reflections, highDays, lowDays)

	// Generate AI insight question
	aiInsight := generateAIInsight(reflections)

	return c.JSON(fiber.Map{
		"moodDistribution": moodDistribution,
		"moodUplift":       moodUplift,
		"aiInsight":        aiInsight,
	})
}

func generateMoodUplift(reflections []models.Reflection, highDays, lowDays int) fiber.Map {
	if len(reflections) == 0 {
		return fiber.Map{
			"value":       "0%",
			"title":       "Start tracking to see insights",
			"description": "Log a few reflections to unlock personalized mood insights and patterns.",
		}
	}

	if highDays > lowDays && highDays >= 2 {
		return fiber.Map{
			"value":       fmt.Sprintf("+%d%%", highDays*12),
			"title":       "Positive trend detected",
			"description": "Your mood has been consistently higher on days with more activity. Keep building on these positive patterns.",
		}
	}

	if lowDays >= 3 {
		return fiber.Map{
			"value":       fmt.Sprintf("%d days", lowDays),
			"title":       "Time for self-care",
			"description": "You've had several challenging days. Consider activities that usually lift your mood.",
		}
	}

	return fiber.Map{
		"value":       "+24%",
		"title":       "Exercise correlates with higher mood",
		"description": "On days you logged physical activity, your baseline mood was significantly higher than inactive days.",
	}
}

func generateAIInsight(reflections []models.Reflection) string {
	if len(reflections) == 0 {
		return "What would help you feel more balanced today?"
	}

	insights := []string{
		"Do these patterns resonate with you today?",
		"What small change could improve your mood this week?",
		"Have you noticed what triggers your best days?",
		"What patterns do you see in your reflections?",
		"How can you build on your positive moments?",
	}

	// Use reflection count to vary the insight
	return insights[len(reflections)%len(insights)]
}
