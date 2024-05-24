package main

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"log"

	"context"
	"os"
	"temperature/internal/infra/cep"
	"temperature/internal/infra/temperature"
	"temperature/internal/infra/weather"
	"temperature/internal/infra/web"
	"temperature/internal/usecase"
)

func main() {
	name := os.Getenv("NAME")
	zipkinUrl := os.Getenv("ZIPKIN_URL")
	_, err := initTracer(name, zipkinUrl)
	if err != nil {
		panic(err)
	}
	tracer := otel.GetTracerProvider().Tracer(name)

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
	webHandler := web.NewWebHandler(getTemperatureUseCase, validateUseCase, tracer)
	err = webHandler.Start(port)
	if err != nil {
		panic(err)
	}
}

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer(name, url string) (func(context.Context) error, error) {
	log.Printf("Starting tracer of %s with url %s", name, url)
	exporter, err := zipkin.New(
		url,
	)
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(bsp),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return tp.Shutdown, nil
}
