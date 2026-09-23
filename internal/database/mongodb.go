package database

import (
	"context"
	"log"
	"time"

	"reflecta/internal/config"

	"go.mongodb.org/mongo-driver/bson"
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

	// Revoked tokens are removed automatically once they would have expired anyway
	_, err = DB.Collection("revoked_tokens").Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.M{"expires_at": 1},
		Options: options.Index().SetExpireAfterSeconds(0),
	})
	if err != nil {
		log.Fatal("❌ Could not create revoked_tokens index:", err)
	}

	log.Println("✅ Successfully connected and pinged MongoDB")
}
