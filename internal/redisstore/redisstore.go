package redisstore

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func New(addr string, password string, db int) *RedisStore {
	return &RedisStore{
		client: redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		}),
	}
}

func (r *RedisStore) SetJSON(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.Set(ctx, key, b, expiration).Err()
}

func (r *RedisStore) GetJSON(ctx context.Context, key string, dest interface{}) (bool, error) {
	res, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal([]byte(res), dest)
}

func (r *RedisStore) HSetJSON(ctx context.Context, key, field string, value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.client.HSet(ctx, key, field, b).Err()
}

func (r *RedisStore) HGetJSON(ctx context.Context, key, field string, dest interface{}) (bool, error) {
	res, err := r.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, json.Unmarshal([]byte(res), dest)
}

func (r *RedisStore) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}

func (r *RedisStore) SAdd(ctx context.Context, key string, members ...string) error {
	return r.client.SAdd(ctx, key, members).Err()
}

func (r *RedisStore) SMembers(ctx context.Context, key string) ([]string, error) {
	return r.client.SMembers(ctx, key).Result()
}
