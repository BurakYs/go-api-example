package authrepository

import (
	"context"
	"errors"
	"time"

	"github.com/BurakYs/go-api-example/internal/database"
	"github.com/BurakYs/go-api-example/internal/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthRepository struct {
	collection *mongo.Collection
	redis      *database.Redis
}

func NewAuthRepository(db *mongo.Database, redis *database.Redis) *AuthRepository {
	return &AuthRepository{
		collection: db.Collection("users"),
		redis:      redis,
	}
}

func (r *AuthRepository) ExistsByUsernameOrEmail(ctx context.Context, username, email string) (bool, error) {
	filter := bson.M{
		"$or": []bson.M{
			{"username": username},
			{"email": email},
		},
	}

	err := r.collection.FindOne(ctx, filter).Err()
	if err == nil {
		return true, nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) {
		return false, nil
	}

	return false, err
}

func (r *AuthRepository) Create(ctx context.Context, user models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	return err
}

func (r *AuthRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *AuthRepository) DeleteByID(ctx context.Context, userID string) error {
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}

	err = r.collection.FindOneAndDelete(ctx, bson.M{
		"_id": objectID,
	}).Err()

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.ErrNotFound
		}

		return err
	}

	sessionsKey := "user_sessions:" + userID
	sessionIDs, err := r.redis.Members(ctx, sessionsKey)
	if err != nil {
		return err
	}

	for _, sessionID := range sessionIDs {
		r.redis.Del(ctx, "session:"+sessionID)
	}

	return nil
}

func (r *AuthRepository) CreateSession(ctx context.Context, userID string) (string, error) {
	sessionID := uuid.NewString()
	if err := r.redis.Set(ctx, "session:"+sessionID, userID, 15*time.Minute); err != nil {
		return "", err
	}

	if err := r.redis.Add(ctx, "user_sessions:"+userID, sessionID); err != nil {
		return "", err
	}

	return sessionID, nil
}

func (r *AuthRepository) DeleteSession(ctx context.Context, sessionID string) error {
	return r.redis.Del(ctx, "session:"+sessionID)
}
