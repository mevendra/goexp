package cep

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
)

type service struct {
	url string
}

func NewCep() Cep {
	return service{
		url: "https://brasilapi.com.br/api/cep/v1/",
	}
}

func (s service) ValidateCep(ctx context.Context, cep string) error {
	_, err := strconv.Atoi(cep)
	if err != nil || len(cep) != 8 {
		return errors.New(InvalidZipCode)
	}
	return nil
}

func (s service) GetCity(ctx context.Context, cep string) (string, error) {
	log.Printf("Executing request for cep api and value '%s'\n", cep)

	url := s.url + cep
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return "", err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if resp.StatusCode != 200 {
		if resp.StatusCode == 404 {
			return "", errors.New(NotFoundZipCode)
		}
		return "", fmt.Errorf("cep api returned status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jsonResp struct {
		City string `json:"city"`
	}
	if err = json.Unmarshal(b, &jsonResp); err != nil {
		return "", err
	}

	return jsonResp.City, nil
}
