package user

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/sixcolors/argon2id"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Register(ctx context.Context, name, email, password string) (*User, error) {
	_, err := s.repo.GetByEmail(ctx, email)
	if err == nil {
		return nil, ErrAlreadyExists
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	hashed, err := s.hashPassword([]byte(password))
	if err != nil {
		return nil, err
	}

	userID, err := s.generateID()
	if err != nil {
		return nil, err
	}

	user := &User{
		ID:        userID,
		Name:      name,
		Email:     email,
		Password:  hashed,
		CreatedAt: time.Now(),
	}

	err = s.repo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email, password string) (*User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	match := s.comparePasswords([]byte(user.Password), []byte(password))
	if !match {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *Service) GetByID(ctx context.Context, id string) (*User, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) generateID() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}

	return id.String(), nil
}

func (s *Service) hashPassword(password []byte) (string, error) {
	b, err := argon2id.GenerateFromPassword(password, nil)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (s *Service) comparePasswords(hashed, password []byte) bool {
	return argon2id.CompareHashAndPassword(hashed, password) == nil
}
