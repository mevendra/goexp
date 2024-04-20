package web

import (
	"net/http"
	"rate_limiter/pkg/entity"
)

type Middleware struct {
	rateLimiter entity.RateLimiter
}

func NewMiddleware(limiter entity.RateLimiter) Middleware {
	return Middleware{rateLimiter: limiter}
}

func (m Middleware) Handle(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.rateLimiter.Allow(r) {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusTooManyRequests)
			return
		}

		h.ServeHTTP(w, r)
	}
}
