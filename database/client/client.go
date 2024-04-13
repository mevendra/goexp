package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	log.Println("Starting client")

	ctx, ctxCancel := context.WithTimeout(context.Background(), time.Millisecond*300)
	defer ctxCancel()
	url := "http://localhost:8080/cotacao"
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		log.Fatal("error creating request", err.Error())
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("error executing request", err.Error())
	}
	if res.StatusCode != http.StatusOK {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		log.Fatal("error reading response body", err.Error())
	}

	var quot Quotation
	err = json.Unmarshal(b, &quot)
	if err != nil {
		log.Fatal("error unmarshalling response", err.Error())
	}

	f, err := os.Create("cotacao.txt")
	if err != nil {
		log.Fatal("error creating file", err.Error())
	}
	defer f.Close()

	_, _ = f.Write([]byte(fmt.Sprintf("Dólar: %.2f", quot.Bid)))
}

type Quotation struct {
	Bid float64 `json:"bid"`
}
