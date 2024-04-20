package entity

import (
	"github.com/stretchr/testify/mock"
	"time"
)

type MemoryMock struct {
	mock.Mock
}

func (m *MemoryMock) Increment(key string) int {
	args := m.Called(key)
	return args.Get(0).(int)
}

func (m *MemoryMock) Expire(key string, ttl time.Duration) {
	_ = m.Called(key, ttl)
}

func (m *MemoryMock) Set(key string, ttl time.Duration) {
	_ = m.Called(key, ttl)
}

func (m *MemoryMock) Get(key string) bool {
	args := m.Called(key)
	return args.Bool(0)
}

func (m *MemoryMock) Del(key string) {
	_ = m.Called(key)
}
