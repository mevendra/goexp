package entity

import (
	"time"
)

type Memory interface {
	Increment(key string) int
	Expire(key string, ttl time.Duration)
	Set(key string, ttl time.Duration)
	Get(key string) bool
	Del(key string)
}
