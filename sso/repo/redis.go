package repo

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStorage struct {
	db *redis.Client
}

func NewSessionStorage(db *redis.Client) *RedisStorage {
	return &RedisStorage{
		db: db,
	}
}

func (r *RedisStorage) Add(ctx context.Context, jwt string, id string, expiresAt time.Duration) error {
	err := r.db.Set(ctx, id, jwt, expiresAt).Err()
	// err := r.db.HSet
	if err != nil {
		return err
	}
	return nil
}

func (r *RedisStorage) Get(ctx context.Context, id string) (string, error) {
	res, err := r.db.Get(ctx, id).Result()
	if err != nil {
		return "", err
	}

	return res, nil

}

func (r *RedisStorage) Delete(ctx context.Context, id string) error {
	return r.db.Del(ctx, id).Err()
}
