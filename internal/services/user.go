package services

import (
	"context"
	"errors"
	"time"

	"github.com/BurakYs/go-api-example/internal/models"
	"github.com/BurakYs/go-api-example/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) GetAll(ctx context.Context, page int) ([]models.PublicUser, error) {
	const pageSize = 10
	skip := int64((page - 1) * pageSize)

	users, err := s.repo.GetAll(ctx, skip, pageSize)
	if err != nil {
		return nil, err
	}

	publicUsers := make([]models.PublicUser, 0, len(users))
	for _, user := range users {
		publicUser := s.convertToPublicUser(user)
		publicUsers = append(publicUsers, publicUser)
	}

	return publicUsers, nil
}

func (s *UserService) GetByID(ctx context.Context, userID string) (*models.PublicUser, error) {
	objectID, err := bson.ObjectIDFromHex(userID)
	if err != nil {
		return nil, models.ErrValidation
	}

	user, err := s.repo.GetByID(ctx, objectID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, models.ErrNotFound
		}

		return nil, err
	}

	publicUser := s.convertToPublicUser(*user)
	return &publicUser, nil
}

func (s *UserService) convertToPublicUser(user models.User) models.PublicUser {
	return models.PublicUser{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Time().Format(time.RFC3339),
	}
}
