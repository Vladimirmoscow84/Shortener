package cache_test

import (
	"context"
	"testing"

	"github.com/Vladimirmoscow84/Shortener.git/internal/storage/cache"
	"github.com/stretchr/testify/require"
	"github.com/wb-go/wbf/redis"
)

// docker run --name redis -p 6379:6379 -d redis:alpine

func TestCache(t *testing.T) {

	client := redis.New(":6379", "", 0)
	require.NotNil(t, client)

	c := cache.NewCache(client)
	ctx := context.Background()

	key := "test-key"
	value := "test-value"

	//Set
	err := c.Set(ctx, key, value)
	require.NoError(t, err, "Set должен успешно записывать значение")

	//Get
	got, err := c.Get(ctx, key)
	require.NoError(t, err, "Get должен успешно читать значение")
	require.Equal(t, value, got, "значения должны совпадать")

	//Exists
	exists, err := c.Exists(ctx, key)
	require.NoError(t, err)
	require.True(t, exists, "ключ должен существовать")

	//отсутствие несуществующего ключа
	missing, err := c.Exists(ctx, "no-such-key")
	require.NoError(t, err)
	require.False(t, missing, "ключа быть не должно")

}
