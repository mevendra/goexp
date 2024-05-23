package cep

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type Mock struct {
	mock.Mock
}

func (m *Mock) ValidateCep(ctx context.Context, cep string) error {
	args := m.Called(ctx, cep)
	return args.Error(0)
}

func (m *Mock) GetCity(ctx context.Context, cep string) (string, error) {
	args := m.Called(ctx, cep)
	return args.String(0), args.Error(1)
}
