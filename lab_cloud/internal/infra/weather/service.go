package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type service struct {
	url    string
	apiKey string // :)
}

func NewWeather() Weather {
	return service{
		url:    "http://api.weatherapi.com/v1",
		apiKey: "b4e406e53a5042b3a2d202757241105",
	}
}

func (s service) GetTemperatureCelsius(ctx context.Context, city string) (float64, error) {
	log.Printf("Executing request for weather api and city '%s'\n", city)

	url := fmt.Sprintf("%s/current.json?key=%s&q=%s", s.url, s.apiKey, city)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return 0, err
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	if resp.StatusCode != 200 {
		return 0, fmt.Errorf("weather api returned status code %d", resp.StatusCode)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var jsonResp struct {
		Current struct {
			TempCelsius float64 `json:"temp_c"`
		} `json:"current"`
	}
	if err = json.Unmarshal(b, &jsonResp); err != nil {
		return 0, err
	}

	return jsonResp.Current.TempCelsius, nil
}
