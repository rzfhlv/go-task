package cache_test

import (
	"context"
	"testing"
	"time"

	"github.com/go-redis/redismock/v9"
	"github.com/rzfhlv/go-task/internal/repository/cache"
	"github.com/stretchr/testify/assert"
)

func TestCacheSet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		mock.ExpectSet("key", int64(1), time.Duration(1*time.Minute)).SetVal("OK")
		cacheRepo := cache.New(client)
		err := cacheRepo.Set(context.Background(), "key", int64(1), time.Duration(1*time.Minute))

		assert.NoError(t, err)
	})
}

func TestCacheGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		mock.ExpectGet("key").SetVal("val")
		cacheRepo := cache.New(client)
		val, err := cacheRepo.Get(context.Background(), "key")

		assert.Equal(t, "val", val)
		assert.Equal(t, nil, err)
	})
}

func TestCacheDel(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		client, mock := redismock.NewClientMock()

		mock.ExpectDel("key").SetVal(int64(1))
		cacheRepo := cache.New(client)
		val, err := cacheRepo.Del(context.Background(), "key")

		assert.Equal(t, int64(1), val)
		assert.Equal(t, nil, err)
	})
}
