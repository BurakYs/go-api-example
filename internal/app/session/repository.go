package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/BurakYs/go-api-example/internal/database"
)

type Repository struct {
	redis *database.Redis
}

func NewRepository(redis *database.Redis) *Repository {
	return &Repository{
		redis: redis,
	}
}

func (r *Repository) Create(ctx context.Context, userID string, expiration time.Duration) (string, error) {
	sessionID := r.generateSessionID()

	err := r.redis.Set(ctx, r.sessionKey(sessionID), userID, expiration)
	if err != nil {
		return "", err
	}

	err = r.redis.SAdd(ctx, r.userSetKey(userID), sessionID)
	if err != nil {
		return "", err
	}

	return sessionID, nil
}

func (r *Repository) Get(ctx context.Context, sessionID string) (string, error) {
	key := r.sessionKey(sessionID)

	val, err := r.redis.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return "", ErrNotFound
		}
		return "", err
	}

	return val, nil
}

func (r *Repository) Delete(ctx context.Context, sessionID string) error {
	key := r.sessionKey(sessionID)

	userID, err := r.redis.Get(ctx, key)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil
		}
		return err
	}

	err = r.redis.Del(ctx, key)
	if err != nil {
		return err
	}

	err = r.redis.SRem(ctx, r.userSetKey(userID), sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteAllForUser(ctx context.Context, userID string) error {
	setKey := r.userSetKey(userID)

	ids, err := r.redis.SMembers(ctx, setKey)
	if err != nil {
		return err
	}

	length := len(ids)
	if length == 0 {
		return r.redis.Del(ctx, setKey)
	}

	keys := make([]string, length+1)
	for i, id := range ids {
		keys[i] = r.sessionKey(id)
	}
	keys[length] = setKey

	return r.redis.Del(ctx, keys...)
}

func (r *Repository) sessionKey(id string) string { return "session:" + id }

func (r *Repository) userSetKey(userID string) string { return "user_sessions:" + userID }

func (r *Repository) generateSessionID() string {
	bytes := make([]byte, 32)
	_, _ = rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
