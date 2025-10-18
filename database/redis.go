package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	client *redis.Client
}

func NewRedis(host, port, password string, db int) (*Redis, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       db,
	})

	err := client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return &Redis{
		client: client,
	}, nil
}

func (r *Redis) Close() error {
	return r.client.Close()
}

func (r *Redis) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *Redis) Get(ctx context.Context, key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Redis) Del(ctx context.Context, keys ...string) error {
	return r.client.Del(ctx, keys...).Err()
}

func (r *Redis) SAdd(ctx context.Context, key string, members ...any) error {
	return r.client.SAdd(ctx, key, members).Err()
}

func (r *Redis) SRem(ctx context.Context, key string, members ...any) error {
	return r.client.SRem(ctx, key, members).Err()
}

func (r *Redis) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}

func (r *Redis) EvalScript(ctx context.Context, script *redis.Script, keys []string, args ...any) ([]int64, error) {
	return script.Run(ctx, r.client, keys, args...).Int64Slice()
}
