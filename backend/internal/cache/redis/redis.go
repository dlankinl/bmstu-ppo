package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"ppo/internal/cache"
	"ppo/internal/config"
	"time"
)

type Cache struct {
	client *redis.Client
}

func NewClient(cfg config.Redis) (client *redis.Client, err error) {
	address := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)
	client = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: cfg.Password,
	})

	status := client.Ping(context.Background())
	_, err = status.Result()
	if err != nil {
		return nil, fmt.Errorf("пинг до redis: %w", err)
	}

	return client, nil
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{client: client}
}

//func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
//	p, err := json.Marshal(value)
//	if err != nil {
//		return fmt.Errorf("сериализация значения: %w", err)
//	}
//
//	err = c.client.Set(ctx, key, p, expiration).Err()
//	if err != nil {
//		return fmt.Errorf("добавление ключа %s в redis кэш: %w", key, err)
//	}
//
//	return nil
//}
//
//func (c *Cache) Get(ctx context.Context, key string, out interface{}) (err error) {
//	err = c.client.Get(ctx, key).Scan(&out)
//	if err != nil {
//		if err == redis.Nil {
//			return cache.ErrNotFound
//		}
//		return fmt.Errorf("получение значения ключа %s из redis кэша: %w", key, err)
//	}
//	fmt.Println("INSIDE GET:", out)
//
//	return nil
//}

func (c *Cache) Get(ctx context.Context, key string, value interface{}) error {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return cache.ErrNotFound
		}
		return fmt.Errorf("получение значения ключа %s из redis кэша %s: %w", key, err)
	}

	err = json.Unmarshal(data, value)
	if err != nil {
		return fmt.Errorf("десериализация значения из redis по ключу %s: %w", key, err)
	}

	return nil
}

func (c *Cache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("сериализация пары ключ-значение (key=%s) в redis: %w", key, err)
	}

	err = c.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return fmt.Errorf("установка значения ключа %s в redis кэш: %w", key, err)
	}

	return nil
}
