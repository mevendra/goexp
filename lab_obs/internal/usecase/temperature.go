package usecase

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"temperature/internal/entity"
	"temperature/internal/infra/cep"
	"temperature/internal/infra/weather"
)

type GetTemperature struct {
	Cep     cep.Cep
	Weather weather.Weather

	tracer trace.Tracer
}

func NewGetTemperatureUseCase(cep cep.Cep, weather weather.Weather, tracer trace.Tracer) *GetTemperature {
	return &GetTemperature{
		Cep:     cep,
		Weather: weather,
		tracer:  tracer,
	}
}

type GetTemperatureInput struct {
	Cep string
}

type GetTemperatureOutput struct {
	City         string  `json:"city"`
	TemperatureC float64 `json:"temp_C"`
	TemperatureF float64 `json:"temp_F"`
	TemperatureK float64 `json:"temp_K"`
}

func (t *GetTemperature) Execute(ctx context.Context, input GetTemperatureInput) (*GetTemperatureOutput, error) {
	ctx, spanCep := t.tracer.Start(ctx, "cep api")
	city, err := t.Cep.GetCity(ctx, input.Cep)
	if err != nil {
		spanCep.End()
		return nil, err
	}
	spanCep.End()

	ctx, spanWeather := t.tracer.Start(ctx, "weather api")
	tmp, err := t.Weather.GetTemperatureCelsius(ctx, city)
	if err != nil {
		spanWeather.End()
		return nil, err
	}
	spanWeather.End()

	output := &GetTemperatureOutput{
		City:         city,
		TemperatureC: tmp,
		TemperatureF: entity.CelsiusToFahrenheit(tmp),
		TemperatureK: entity.CelsiusToKelvin(tmp),
	}
	return output, nil
}
