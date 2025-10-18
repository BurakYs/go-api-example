package user

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/BurakYs/go-api-example/database"
)

type Repository struct {
	collection *mongo.Collection
}

func NewRepository(db *database.DB) *Repository {
	return &Repository{
		collection: db.GetCollection("users"),
	}
}

func (r *Repository) Create(ctx context.Context, user *User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil && mongo.IsDuplicateKeyError(err) {
		return ErrAlreadyExists
	}

	return err
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (*User, error) {
	return r.getByFilter(ctx, bson.M{"email": email})
}

func (r *Repository) GetByID(ctx context.Context, id string) (*User, error) {
	return r.getByFilter(ctx, bson.M{"_id": id})
}

func (r *Repository) CreateIndexes(ctx context.Context) error {
	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{
				{Key: "email", Value: 1},
			},
			Options: options.Index().SetUnique(true).SetName("email_index"),
		},
	}

	_, err := r.collection.Indexes().CreateMany(ctx, indexes)
	return err
}

func (r *Repository) getByFilter(ctx context.Context, filter any) (*User, error) {
	var user User

	err := r.collection.FindOne(ctx, filter).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
