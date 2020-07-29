package storage

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

// Redis is a wrapper over the redis client.
type Redis struct {
	client *redis.Client
}

// NewRedis creates a new instance of Redis storage.
func NewRedis(options *redis.Options) *Redis {
	return &Redis{
		client: redis.NewClient(options),
	}
}

// Get returns the pre-cached value for the given key.
func (r *Redis) Get(key string) (string, error) {
	val, err := r.client.Get(context.TODO(), key).Result()
	return val, err
}

// Set sets the value for the specified key in the cache.
func (r *Redis) Set(key string, value []byte, expiration time.Duration) error {
	return r.client.Set(context.TODO(), key, value, expiration).Err()
}

// Delete removes the entry for the specified key.
func (r *Redis) Delete(key string) error {
	return r.client.Del(context.TODO(), key).Err()
}
