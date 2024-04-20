package entity

import (
	"errors"
	"net/http"
	"rate_limiter/config"
	"strings"
	"time"
)

type RateLimiter struct {
	token      bool
	tokenLimit int
	tokenBlock time.Duration
	tokenFrame time.Duration
	ip         bool
	ipLimit    int
	ipBlock    time.Duration
	ipFrame    time.Duration
	memory     Memory
}

func NewRateLimiter(conf *config.Config, memory Memory) RateLimiter {
	isToken := conf.TokenLimit > 0
	if isToken {
		if conf.TokenBlockTime <= 0 || conf.TokenFrameTime <= 0 {
			panic(errors.New("invalid token block or frame time"))
		}
	}

	isIp := conf.IpLimit > 0
	if isIp {
		if conf.IpBlockTime <= 0 || conf.IpFrameTime <= 0 {
			panic(errors.New("invalid ip block or frame time"))
		}
	}

	return RateLimiter{
		token:      isToken,
		tokenLimit: conf.TokenLimit,
		tokenBlock: conf.TokenBlockTime,
		tokenFrame: conf.TokenFrameTime,

		ip:      isIp,
		ipLimit: conf.IpLimit,
		ipBlock: conf.IpBlockTime,
		ipFrame: conf.IpFrameTime,

		memory: memory,
	}
}

func (l *RateLimiter) Allow(r *http.Request) bool {
	if l.token {
		key := r.Header.Get("API_KEY")
		return l.allow(key, l.tokenLimit, l.tokenBlock, l.tokenFrame)
	}

	if l.ip {
		key := strings.Split(r.RemoteAddr, ":")[0]
		return l.allow(key, l.ipLimit, l.ipBlock, l.ipBlock)
	}

	return true
}

func (l *RateLimiter) allow(key string, limit int, blockTime time.Duration, frameTime time.Duration) bool {
	blockedKey := l.getBlockedKey(key)
	if l.memory.Get(blockedKey) {
		return false
	}

	current := l.memory.Increment(key)
	if current <= 1 {
		l.memory.Expire(key, frameTime)
	}

	if current > limit {
		l.memory.Set(blockedKey, blockTime)
		return false
	}

	return true
}

func (l *RateLimiter) getBlockedKey(key string) string {
	return "blocked_" + key
}
