package usecase

import (
	"context"
	"temperature/internal/infra/cep"
	"temperature/internal/infra/temperature"
)

type ValidateCep struct {
	Cep cep.Cep
	Tmp temperature.Temperature
}

func NewValidateCepUseCase(cep cep.Cep, tmp temperature.Temperature) *ValidateCep {
	return &ValidateCep{Cep: cep, Tmp: tmp}
}

type ValidateCepInput struct {
	Cep string
}

type ValidateCepOutput struct {
	City         string  `json:"city"`
	TemperatureC float64 `json:"temp_C"`
	TemperatureF float64 `json:"temp_F"`
	TemperatureK float64 `json:"temp_K"`
}

func (v *ValidateCep) Execute(ctx context.Context, input ValidateCepInput) (*ValidateCepOutput, error) {
	err := v.Cep.ValidateCep(ctx, input.Cep)
	if err != nil {
		return nil, err
	}

	o, err := v.Tmp.GetTemperature(ctx, input.Cep)
	if err != nil {
		return nil, err
	}

	return &ValidateCepOutput{
		City:         o.City,
		TemperatureC: o.C,
		TemperatureF: o.F,
		TemperatureK: o.K,
	}, nil
}
