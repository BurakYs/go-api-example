package session

import (
	"context"
	"time"
)

type Service struct {
	repo       *Repository
	expiration time.Duration
}

func NewService(repo *Repository, expiration time.Duration) *Service {
	return &Service{
		repo:       repo,
		expiration: expiration,
	}
}

func (s *Service) Create(ctx context.Context, userID string) (string, error) {
	return s.repo.Create(ctx, userID, s.expiration)
}

func (s *Service) GetUserID(ctx context.Context, sessionID string) (string, error) {
	return s.repo.Get(ctx, sessionID)
}

func (s *Service) Delete(ctx context.Context, sessionID string) error {
	return s.repo.Delete(ctx, sessionID)
}

func (s *Service) RevokeAllForUser(ctx context.Context, userID string) error {
	return s.repo.DeleteAllForUser(ctx, userID)
}
