package main

import (
	"os"
	"temperature/internal/infra/cep"
	"temperature/internal/infra/temperature"
	"temperature/internal/infra/weather"
	"temperature/internal/infra/web"
	"temperature/internal/usecase"
)

func main() {
	weatherService := weather.NewWeather()
	cepService := cep.NewCep()
	temperatureUri := os.Getenv("TEMPERATURE_URI")
	if temperatureUri == "" {
		temperatureUri = "http://localhost:8081"
	}
	temperatureService := temperature.NewTemperature(temperatureUri)

	getTemperatureUseCase := usecase.NewGetTemperatureUseCase(cepService, weatherService)
	validateUseCase := usecase.NewValidateCepUseCase(cepService, temperatureService)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	webHandler := web.NewWebHandler(getTemperatureUseCase, validateUseCase)
	err := webHandler.Start(port)
	if err != nil {
		panic(err)
	}
}
