package main

import (
	"temperature/internal/infra/cep"
	"temperature/internal/infra/weather"
	"temperature/internal/infra/web"
	"temperature/internal/usecase"
)

func main() {
	weatherService := weather.NewWeather()
	cepService := cep.NewCep()
	getTemperatureUseCase := usecase.NewGetTemperatureUseCase(cepService, weatherService)

	port := "8080"
	webHandler := web.NewWebHandler(getTemperatureUseCase)
	err := webHandler.Start(port)
	if err != nil {
		panic(err)
	}
}
