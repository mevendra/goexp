package web

import (
	"encoding/json"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"io"
	"log"
	"net/http"
	"temperature/internal/infra/cep"
	"temperature/internal/usecase"
)

type Handler struct {
	temperature *usecase.GetTemperature
	validate    *usecase.ValidateCep

	tracer trace.Tracer
}

func NewWebHandler(temperature *usecase.GetTemperature, validate *usecase.ValidateCep, tracer trace.Tracer) Handler {
	return Handler{
		temperature: temperature,
		validate:    validate,
		tracer:      tracer,
	}
}

func (h Handler) serveHTTPTemperature(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.tracer.Start(ctx, "temperature")
	defer span.End()

	log.Println("Serving HTTP temperature request")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cepInput := r.URL.Query().Get("cep")

	input := usecase.GetTemperatureInput{Cep: cepInput}
	output, err := h.temperature.Execute(ctx, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	b, err := json.Marshal(output)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h Handler) serveHTTPValidate(w http.ResponseWriter, r *http.Request) {
	carrier := propagation.HeaderCarrier(r.Header)
	ctx := r.Context()
	ctx = otel.GetTextMapPropagator().Extract(ctx, carrier)
	ctx, span := h.tracer.Start(ctx, "validate")
	defer span.End()

	log.Println("Serving HTTP validate request")
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	if r.Body == nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		h.handleError(w, err)
		return
	}

	rInput := struct {
		Cep string `json:"cep"`
	}{}
	err = json.Unmarshal(b, &rInput)
	if err != nil {
		h.handleError(w, err)
		return
	}

	input := usecase.ValidateCepInput{Cep: rInput.Cep}
	output, err := h.validate.Execute(ctx, input)
	if err != nil {
		h.handleError(w, err)
		return
	}

	b, err = json.Marshal(output)
	if err != nil {
		h.handleError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h Handler) handleError(w http.ResponseWriter, err error) {
	msg := err.Error()
	switch msg {
	case cep.InvalidZipCode:
		http.Error(w, msg, http.StatusUnprocessableEntity)
	case cep.NotFoundZipCode:
		http.Error(w, msg, http.StatusNotFound)
	default:
		http.Error(w, msg, http.StatusInternalServerError)
	}
	return
}

func (h Handler) Start(port string) error {
	log.Printf("Starting web server on port %s", port)

	http.HandleFunc("/temperature", h.serveHTTPTemperature)
	http.HandleFunc("/validate", h.serveHTTPValidate)
	return http.ListenAndServe(":"+port, nil)
}
