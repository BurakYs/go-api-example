package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type DB struct {
	client   *mongo.Client
	database *mongo.Database
}

func NewDB(uri, dbName string) (*DB, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err
	}

	return &DB{
		client:   client,
		database: client.Database(dbName),
	}, nil
}

func (d *DB) GetCollection(name string) *mongo.Collection {
	return d.database.Collection(name)
}

func (d *DB) Disconnect(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}
