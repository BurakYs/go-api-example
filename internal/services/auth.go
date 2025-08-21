package services

import (
	"context"
	"errors"
	"time"

	"github.com/BurakYs/go-api-example/internal/models"
	"github.com/BurakYs/go-api-example/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo   *repository.AuthRepository
	domain string
}

func NewAuthService(repo *repository.AuthRepository, domain string) *AuthService {
	return &AuthService{
		repo:   repo,
		domain: domain,
	}
}

func (s *AuthService) Register(ctx context.Context, body *models.RegisterUserBody) (models.PublicUser, string, error) {
	exists, err := s.repo.ExistsByUsernameOrEmail(ctx, body.Username, body.Email)
	if err != nil {
		return models.PublicUser{}, "", err
	}

	if exists {
		return models.PublicUser{}, "", models.ErrConflict
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return models.PublicUser{}, "", err
	}

	user := models.User{
		ID:        bson.NewObjectID(),
		Username:  body.Username,
		Email:     body.Email,
		Password:  string(hashedPassword),
		CreatedAt: bson.DateTime(time.Now().UTC().UnixMilli()),
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return models.PublicUser{}, "", err
	}

	sessionID, err := s.repo.CreateSession(ctx, user.ID.Hex())
	if err != nil {
		return models.PublicUser{}, "", err
	}

	return models.PublicUser{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Time().Format(time.RFC3339),
	}, sessionID, nil
}

func (s *AuthService) Login(ctx context.Context, body *models.LoginUserBody) (models.PublicUser, string, error) {
	user, err := s.repo.GetByEmail(ctx, body.Email)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return models.PublicUser{}, "", models.ErrNotFound
		}

		return models.PublicUser{}, "", err
	}

	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)) != nil {
		return models.PublicUser{}, "", models.ErrForbidden
	}

	sessionID, err := s.repo.CreateSession(ctx, user.ID.Hex())
	if err != nil {
		return models.PublicUser{}, "", err
	}

	return models.PublicUser{
		ID:        user.ID.Hex(),
		Username:  user.Username,
		CreatedAt: user.CreatedAt.Time().Format(time.RFC3339),
	}, sessionID, nil
}

func (s *AuthService) Logout(ctx context.Context, sessionID string) error {
	return s.repo.DeleteSession(ctx, sessionID)
}

func (s *AuthService) DeleteAccount(ctx context.Context, userID string) error {
	err := s.repo.DeleteByID(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}

func (s *AuthService) Domain() string {
	return s.domain
}
