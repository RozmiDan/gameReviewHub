package redis_build

import (
	"context"
	"time"

	"github.com/RozmiDan/gameReviewHub/internal/entity"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisCache struct {
	client *redis.Client
	logger *zap.Logger
	ttl    int
}

func NewRedisClient(redis_address, passwrd string, db, ttl int, logger *zap.Logger) *RedisCache {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_address,
		Password: passwrd,
		DB:       db,
	})
	logger = logger.With(zap.String("component", "Redis"))
	return &RedisCache{client: client, logger: logger, ttl: ttl}
}

func (r *RedisCache) Get(ctx context.Context, key string) (string, error) {
	newCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	s, err := r.client.Get(newCtx, key).Result()
	if err == redis.Nil {
		return "", entity.ErrCacheMiss
	}
	if err != nil {
		r.logger.Info("cant get values from redis", zap.Error(err))
		return "", err
	}
	return s, nil
}

func (r *RedisCache) Set(ctx context.Context, key, value string) error {
	newCtx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	err := r.client.Set(newCtx, key, value, time.Duration(r.ttl)*time.Second).Err()
	if err != nil {
		r.logger.Error("cant set values in redis", zap.Error(err))
		return err
	}
	return nil
}
