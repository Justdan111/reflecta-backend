package database

import (
	"context"
	"log"
	"time"

	"reflecta/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database

func ConnectMongo() {
	uri := config.GetEnv("MONGO_URI")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal("❌ Error connecting to MongoDB:", err)
	}

	// Ping the database to confirm connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("❌ Could not ping MongoDB:", err)
	}

	DB = client.Database("reflecta")

	log.Println("✅ Successfully connected and pinged MongoDB")
}
