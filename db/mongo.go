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
var Collections = struct {
	Users *mongo.Collection
}{
	Users: nil,
}

func SetupMongo() {
	client, err := mongo.Connect(options.Client().ApplyURI(config.AppConfig.MongoURI))
	if err != nil {
		log.Fatalln("Failed to connect to MongoDB:", err)
	}

	Mongo = client

	Collections.Users = getCollection("users")

	_, err = Collections.Users.Indexes().CreateMany(context.Background(), []mongo.IndexModel{
		{
			Keys:    bson.M{"email": 1},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys:    bson.M{"username": 1},
			Options: options.Index().SetUnique(true),
		},
	})

	if err != nil {
		log.Fatalln("Error creating indexes:", err)
	}
}

func DisconnectMongo() error {
	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := Mongo.Disconnect(timeoutCtx); err != nil {
		log.Println("Error disconnecting from MongoDB:", err)
		return err
	}

	return nil
}

func getCollection(name string) *mongo.Collection {
	return Mongo.Database(config.AppConfig.MongoDBName).Collection(name)
}
