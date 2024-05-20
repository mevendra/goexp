package weather

import "context"

type Weather interface {
	GetTemperatureCelsius(ctx context.Context, city string) (float64, error)
}
