package cache

import (
	"context"
	"errors"
	"time"
)

var ErrNotFound = errors.New("в redis-кэше нет такого ключа")

type Cache interface {
	Get(ctx context.Context, key string, out interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}
