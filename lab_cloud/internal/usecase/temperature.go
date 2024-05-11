package usecase

import (
	"context"
	"temperature/internal/entity"
	"temperature/internal/infra/cep"
	"temperature/internal/infra/weather"
)

type GetTemperature struct {
	Cep     cep.Cep
	Weather weather.Weather
}

func NewGetTemperatureUseCase(cep cep.Cep, weather weather.Weather) *GetTemperature {
	return &GetTemperature{
		Cep:     cep,
		Weather: weather,
	}
}

type GetTemperatureInput struct {
	Cep string
}

type GetTemperatureOutput struct {
	TemperatureC float64 `json:"temp_C"`
	TemperatureF float64 `json:"temp_F"`
	TemperatureK float64 `json:"temp_K"`
}

func (t *GetTemperature) Execute(ctx context.Context, input GetTemperatureInput) (*GetTemperatureOutput, error) {
	c := input.Cep
	if err := t.Cep.ValidateCep(ctx, c); err != nil {
		return nil, err
	}

	city, err := t.Cep.GetCity(ctx, c)
	if err != nil {
		return nil, err
	}

	tmp, err := t.Weather.GetTemperatureCelsius(ctx, city)
	if err != nil {
		return nil, err
	}

	output := &GetTemperatureOutput{
		TemperatureC: tmp,
		TemperatureF: entity.CelsiusToFahrenheit(tmp),
		TemperatureK: entity.CelsiusToKelvin(tmp),
	}
	return output, nil
}
