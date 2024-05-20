package temperature

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	netUrl "net/url"
)

type service struct {
	uri string
}

func NewTemperature(uri string) Temperature {
	log.Printf("Starting temperature service on uri %s", uri)
	return &service{uri: uri}
}

func (s service) GetTemperature(ctx context.Context, cep string) (Output, error) {
	log.Printf("Executing request for temperature api for cep %s", cep)
	cep = netUrl.QueryEscape(cep)
	url := fmt.Sprintf("%s/temperature?cep=%s", s.uri, cep)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Output{}, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return Output{}, err
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return Output{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Output{}, errors.New(string(b))
	}

	var jsonResp struct {
		City  string  `json:"city"`
		TempC float64 `json:"temp_C"`
		TempF float64 `json:"temp_F"`
		TempK float64 `json:"temp_K"`
	}
	if err = json.Unmarshal(b, &jsonResp); err != nil {
		return Output{}, err
	}
	return Output{
		City: jsonResp.City,
		C:    jsonResp.TempC,
		F:    jsonResp.TempF,
		K:    jsonResp.TempK,
	}, nil
}
