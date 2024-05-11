package web

import (
	"encoding/json"
	"log"
	"net/http"
	"temperature/internal/infra/cep"
	"temperature/internal/usecase"
)

type Handler struct {
	temperature *usecase.GetTemperature
}

func NewWebHandler(temperature *usecase.GetTemperature) Handler {
	return Handler{
		temperature: temperature,
	}
}

func (h Handler) serveHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving HTTP request")
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	ctx := r.Context()
	cepInput := r.URL.Query().Get("cep")

	input := usecase.GetTemperatureInput{Cep: cepInput}
	output, err := h.temperature.Execute(ctx, input)
	if err != nil {
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

	b, err := json.Marshal(output)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(b)
}

func (h Handler) Start(port string) error {
	log.Printf("Starting web server on port %s", port)

	http.HandleFunc("/", h.serveHTTP)
	return http.ListenAndServe(":"+port, nil)
}
