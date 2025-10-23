package cache

import (
	"context"

	"github.com/wb-go/wbf/redis"
)

type Cache struct {
	client *redis.Client
}

// NewCache создает новый кэш
func NewCache(client *redis.Client) *Cache {
	return &Cache{
		client: client,
	}
}

// Set устанавливает значение в кэш по ключу или перезаписывает
func (c *Cache) Set(ctx context.Context, key string, value any) error {
	return c.client.Set(ctx, key, value)
}

// Get возвращает значение из кэш по ключу
func (c *Cache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key)
}

// Exists проверяет наличие записей в кэш по ключу
func (c *Cache) Exists(ctx context.Context, key string) (bool, error) {
	count, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
