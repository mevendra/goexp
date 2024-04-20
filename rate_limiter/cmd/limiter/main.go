package main

import (
	"log"
	"net/http"
	"rate_limiter/config"
	"rate_limiter/pkg/entity"
	"rate_limiter/pkg/infra/redis"
	"rate_limiter/pkg/infra/web"
	"strconv"
)

func main() {
	conf := config.LoadConfig(".")
	log.Printf("Server configurations: %s\n", conf.String())

	memory := redis.NewMemory(conf)
	rateLimiter := entity.NewRateLimiter(conf, memory)

	httpMiddleware := web.NewMiddleware(rateLimiter)
	http.Handle("/", httpMiddleware.Handle(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))

	addr := ":" + strconv.Itoa(conf.Port)
	log.Printf("Starting limiter test server on port %s\n", addr)
	_ = http.ListenAndServe(addr, nil)
}
