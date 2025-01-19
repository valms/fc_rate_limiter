package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

type RedisRepository struct {
	ctx         context.Context
	redisClient *redis.Client
}

func NewRedisOperations(ctx context.Context, redisClient *redis.Client) (*RedisRepository, error) {
	if ctx == nil {
		return nil, fmt.Errorf("context cannot be nil")
	}
	if redisClient == nil {
		return nil, fmt.Errorf("redis cache cannot be nil")
	}

	return &RedisRepository{
		ctx:         ctx,
		redisClient: redisClient,
	}, nil
}

// GetCount retorna o valor atual do contador
func (r *RedisRepository) GetCounter(key string) (int, error) {
	if key == "" {
		return 0, fmt.Errorf("key cannot be empty")
	}

	result, err := r.redisClient.Get(r.ctx, key).Int()
	if err == redis.Nil {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("failed to get count for key %s: %w", key, err)
	}

	return result, nil
}

// Increment incrementa o valor da chave em 1
func (r *RedisRepository) IncrementCounter(key string) (int, error) {
	if key == "" {
		return 0, fmt.Errorf("key cannot be empty")
	}

	result, err := r.redisClient.IncrBy(r.ctx, key, 1).Result()
	if err != nil {
		return 0, fmt.Errorf("failed to increment key %s: %w", key, err)
	}
	return int(result), nil
}

// Expire define um tempo de expiração para a chave
func (r *RedisRepository) SetExpiration(key string, seconds int) error {
	if key == "" {
		return fmt.Errorf("key cannot be empty")
	}

	if seconds <= 0 {
		return fmt.Errorf("expiration time must be positive")
	}

	duration := time.Duration(seconds) * time.Second
	success, err := r.redisClient.Expire(r.ctx, key, duration).Result()
	if err != nil {
		return fmt.Errorf("failed to set expiry for key %s: %w", key, err)
	}
	if !success {
		// Se a chave não existir, cria ela com o valor 1
		err = r.redisClient.Set(r.ctx, key, 1, duration).Err()
		if err != nil {
			return fmt.Errorf("failed to create and set expiry for key %s: %w", key, err)
		}
	}
	return nil
}

// Exists verifica se uma chave existe
func (r *RedisRepository) KeyExists(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key cannot be empty")
	}

	result, err := r.redisClient.Exists(r.ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check existence of key %s: %w", key, err)
	}

	return result > 0, nil
}

// Close fecha a conexão com o Redis
func (r *RedisRepository) Close() error {
	if r.redisClient == nil {
		return nil
	}

	if err := r.redisClient.Close(); err != nil {
		return fmt.Errorf("failed to close redis cache: %w", err)
	}
	return nil
}
