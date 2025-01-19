package cache

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(addr, password string) (*redis.Client, error) {
	fmt.Printf("Connecting to redis at %s\n", addr)
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	ctx := context.Background()

	if err := client.Ping(ctx).Err(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	fmt.Printf("Connected to redis at %s\n", addr)

	return client, nil
}
