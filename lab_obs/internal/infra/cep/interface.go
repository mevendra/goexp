package cep

import (
	"context"
)

const (
	InvalidZipCode  = "invalid zipcode"
	NotFoundZipCode = "can not find zipcode"
)

type Cep interface {
	ValidateCep(ctx context.Context, cep string) error
	GetCity(ctx context.Context, cep string) (string, error)
}
