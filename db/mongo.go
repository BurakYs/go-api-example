package db

import (
	"context"
	"log"
	"time"

	"github.com/BurakYs/GoAPIExample/config"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Mongo *mongo.Client

func SetupMongo() {
	client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		log.Fatal("MongoDB connection error:", err)
	}

	Mongo = client

	_, err = GetCollection("users").Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.D{{Key: "username", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
	})

	if err != nil {
		log.Fatalf("Error creating indexes: %v", err)
	}
}

func DisconnectMongo() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Mongo.Disconnect(timeoutCtx); err != nil {
		log.Printf("Error disconnecting from MongoDB: %v", err)
		return err
	}

	log.Println("MongoDB disconnected")
	return nil
}

func GetCollection(collectionName string) *mongo.Collection {
	return Mongo.Database(config.AppConfig.MongoDBName).Collection(collectionName)
}
