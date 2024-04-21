package entity

import (
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"testing"
	"time"
)

func BenchmarkRateLimiter_AllowToken(b *testing.B) {
	mem := new(MemoryMock)
	mem.On("Increment", mock.Anything).Return(1)
	mem.On("Expire", mock.Anything, mock.Anything)
	mem.On("Get", mock.Anything).Return(false)
	rateLimiter := RateLimiter{
		token:      true,
		tokenLimit: 100,
		tokenBlock: time.Millisecond,
		tokenFrame: time.Millisecond,
		memory:     mem,
	}

	wg := sync.WaitGroup{}
	wg.Add(b.N * b.N)
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("API_KEY", strconv.Itoa(i))
		for j := 0; j < b.N; j++ {
			go func() {
				rateLimiter.Allow(req)
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func BenchmarkRateLimiter_AllowIp(b *testing.B) {
	mem := new(MemoryMock)
	mem.On("Increment", mock.Anything).Return(1)
	mem.On("Expire", mock.Anything, mock.Anything)
	mem.On("Get", mock.Anything).Return(false)
	rateLimiter := RateLimiter{
		ip:      true,
		ipLimit: 100,
		ipBlock: time.Millisecond,
		ipFrame: time.Millisecond,
		memory:  mem,
	}

	wg := sync.WaitGroup{}
	wg.Add(b.N * b.N)
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.RemoteAddr = "127.0.0.1:8080"
		for j := 0; j < b.N; j++ {
			go func() {
				rateLimiter.Allow(req)
				wg.Done()
			}()
		}
	}

	wg.Wait()
}

func BenchmarkRateLimiter_AllowBlocked(b *testing.B) {
	mem := new(MemoryMock)
	mem.On("Get", mock.Anything).Return(true)
	rateLimiter := RateLimiter{
		token:      true,
		tokenLimit: 100,
		tokenBlock: time.Millisecond,
		tokenFrame: time.Millisecond,
		memory:     mem,
	}

	wg := sync.WaitGroup{}
	wg.Add(b.N * b.N)
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Add("API_KEY", strconv.Itoa(i))
		for j := 0; j < b.N; j++ {
			go func() {
				rateLimiter.Allow(req)
				wg.Done()
			}()
		}
	}

	wg.Wait()
}
