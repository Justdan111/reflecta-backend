package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"refecta/internal/config"
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

	DB = client.Database("reflecta")

	log.Println("✅ Connected to MongoDB")
}