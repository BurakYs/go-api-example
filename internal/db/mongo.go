package db

import (
	"context"
	"log"

	"github.com/BurakYs/GoAPIExample/internal/config"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var mongoClient *mongo.Client

func SetupMongo() {
	client, err := mongo.Connect(options.Client().ApplyURI(config.App.MongoURI))
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
	}

	mongoClient = client
}

func DisconnectMongo() {
	mongoClient.Disconnect(context.Background())
}

func GetCollection(name string) *mongo.Collection {
	return mongoClient.Database(config.App.MongoDBName).Collection(name)
}
