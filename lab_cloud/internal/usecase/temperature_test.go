package usecase

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"temperature/internal/infra/cep"
	"temperature/internal/infra/weather"
	"testing"
)

func TestGetTemperature_Execute(t *testing.T) {
	t.Run("Valid input", func(t *testing.T) {
		expected := &GetTemperatureOutput{
			TemperatureC: 24,
			TemperatureF: 75.2,
			TemperatureK: 297,
		}
		input := GetTemperatureInput{
			Cep: "99500000",
		}

		cepMock := new(cep.Mock)
		cepMock.On("ValidateCep", mock.Anything, "99500000").Return(nil)
		cepMock.On("GetCity", mock.Anything, input.Cep).Return("City", nil)

		weatherMock := new(weather.Mock)
		weatherMock.On("GetTemperatureCelsius", mock.Anything, "City").Return(24.0, nil)
		uCase := &GetTemperature{
			Cep:     cepMock,
			Weather: weatherMock,
		}

		actual, err := uCase.Execute(context.Background(), input)
		assert.Nil(t, err)
		assert.Equal(t, expected, actual)
	})

	t.Run("Invalid Cep", func(t *testing.T) {
		expectedErr := cep.InvalidZipCode
		input := GetTemperatureInput{
			Cep: "99500000",
		}

		cepMock := new(cep.Mock)
		cepMock.On("ValidateCep", mock.Anything, "99500000").Return(errors.New(cep.InvalidZipCode))
		uCase := &GetTemperature{
			Cep: cepMock,
		}

		_, err := uCase.Execute(context.Background(), input)
		if assert.NotNil(t, err) {
			assert.Equal(t, expectedErr, err.Error())
		}
	})

	t.Run("Invalid City", func(t *testing.T) {
		expectedErr := cep.NotFoundZipCode
		input := GetTemperatureInput{
			Cep: "99500000",
		}

		cepMock := new(cep.Mock)
		cepMock.On("ValidateCep", mock.Anything, "99500000").Return(nil)
		cepMock.On("GetCity", mock.Anything, input.Cep).Return("", errors.New(cep.NotFoundZipCode))

		weatherMock := new(weather.Mock)
		weatherMock.On("GetTemperatureCelsius", mock.Anything, "City").Return(24.0, nil)
		uCase := &GetTemperature{
			Cep:     cepMock,
			Weather: weatherMock,
		}

		_, err := uCase.Execute(context.Background(), input)
		if assert.NotNil(t, err) {
			assert.Equal(t, expectedErr, err.Error())
		}
	})

	t.Run("Invalid Weather", func(t *testing.T) {
		expectedErr := "invalid weather"
		input := GetTemperatureInput{
			Cep: "99500000",
		}

		cepMock := new(cep.Mock)
		cepMock.On("ValidateCep", mock.Anything, "99500000").Return(nil)
		cepMock.On("GetCity", mock.Anything, input.Cep).Return("City", nil)

		weatherMock := new(weather.Mock)
		weatherMock.On("GetTemperatureCelsius", mock.Anything, "City").Return(0.0, errors.New("invalid weather"))
		uCase := &GetTemperature{
			Cep:     cepMock,
			Weather: weatherMock,
		}

		_, err := uCase.Execute(context.Background(), input)
		if assert.NotNil(t, err) {
			assert.Equal(t, expectedErr, err.Error())
		}
	})
}
