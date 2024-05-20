package temperature

import "context"

type Temperature interface {
	GetTemperature(ctx context.Context, cep string) (Output, error)
}

type Output struct {
	City string
	C    float64
	F    float64
	K    float64
}
