package redis

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log/slog"
	"time"
)

type Cache struct {
	client *redis.Client
	logger *slog.Logger
	config CacheConfig
}

type CacheConfig struct {
	TTL      time.Duration `koanf:"ttl"`
	WriteTTL time.Duration `koanf:"write_ttl"`
}

const tag = "cache_manager"

func NewRedisCache(client *redis.Client, logger *slog.Logger, config CacheConfig) *Cache {
	return &Cache{
		client: client,
		logger: logger,
		config: config,
	}
}

func (r *Cache) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := r.client.Get(ctx, key).Bytes()
	if err != nil {
		//todo add metrics --> cache miss
		r.logger.Info("cache miss", "tag", tag, "key", key)
		return nil, err
	}

	r.logger.Info("Cache hit", "tag", tag, "key", key)
	//todo add metrics --> cache hit
	return data, nil
}

func (r *Cache) Set(ctx context.Context, key string, value interface{}) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.config.WriteTTL)
	defer cancel()
	err := r.client.Set(ctx, key, value, r.config.TTL).Err()
	if err != nil {
		// todo add metrics --> setting in cache failed
		r.logger.Error("error setting in cache", "tag", tag, "key", key, "err", err)
	} else {
		// todo add metrics --> setting in cache succeed
		r.logger.Info("set in cache succeed", "tag", tag, "key", key)
	}
	return err
}

func (r *Cache) Invalidate(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		// todo add metrics --> error invalidating cache
		r.logger.Error("invalidating cache", "tag", tag, "key", key, "err", err)
	} else {
		// todo add metrics --> invalidating cache succeed
		r.logger.Info("invalidating cache", "tag", tag, "key", key)
	}
	return err
}
