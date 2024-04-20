package redis

import (
	"context"
	"errors"
	r "github.com/go-redis/redis/v8"
	"rate_limiter/config"
	"time"
)

type Memory struct {
	client *r.Client
}

func NewMemory(conf *config.Config) *Memory {
	client := r.NewClient(&r.Options{
		Addr:     conf.RedisAddr,
		Username: conf.RedisUsername,
		Password: conf.RedisPassword,
	})

	status := client.Ping(context.Background())
	if _, err := status.Result(); err != nil {
		panic(err)
	}

	return &Memory{
		client: client,
	}
}

func (m Memory) Increment(key string) int {
	newValue, err := m.client.Incr(context.Background(), key).Result()
	if err != nil {
		return 0
	}
	return int(newValue)
}

func (m Memory) Expire(key string, ttl time.Duration) {
	_ = m.client.Expire(context.Background(), key, ttl)
}

func (m Memory) Set(key string, ttl time.Duration) {
	m.client.Set(context.Background(), key, key, ttl)
}

func (m Memory) Get(key string) bool {
	_, err := m.client.Get(context.Background(), key).Result()
	if err != nil {
		return !errors.Is(err, r.Nil)
	}
	return true
}

func (m Memory) Del(key string) {
	_ = m.client.Del(context.Background(), key)
}
