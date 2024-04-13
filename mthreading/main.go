package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

func main() {
	log.Println("Starting mthreading")
	cep := Cep{}
	cep.AddUnformattedUrl("Brasil API", "https://brasilapi.com.br/api/cep/v1/%d")
	cep.AddUnformattedUrl("Via CEP", "http://viacep.com.br/ws/%d/json/")

	cepHandler := &CepHandler{C: cep}

	http.Handle("/", cepHandler)
	_ = http.ListenAndServe(":8080", nil)

	log.Println("End of mthreading")
}

type CepHandler struct {
	C Cep
}

func (h *CepHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var cep int
	var err error
	cepKey := "cep"
	cepStr := r.URL.Query().Get(cepKey)
	if cep, err = strconv.Atoi(cepStr); err != nil {
		cepStr = r.Header.Get(cepKey)
		if cep, err = strconv.Atoi(cepStr); err != nil {
			path := r.URL.Path
			pathSplit := strings.Split(path, "/")
			if len(pathSplit) <= 0 {
				http.Error(w, "invalid path", http.StatusBadRequest)
				return
			}

			cepStr = pathSplit[len(pathSplit)-1]
			if cep, err = strconv.Atoi(cepStr); err != nil {
				http.Error(w, "invalid cep parameter", http.StatusBadRequest)
				return
			}
		}
	}

	res, err := h.C.GetCep(cep)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(200)
	_, _ = w.Write(res)
}

type Cep struct {
	Names           []string
	UnformattedUrls []string
}

func (c *Cep) AddUnformattedUrl(name, url string) {
	c.Names = append(c.Names, name)
	c.UnformattedUrls = append(c.UnformattedUrls, url)
}

func (c *Cep) GetCep(n int) ([]byte, error) {
	wg := &sync.WaitGroup{}
	ch := make(chan []byte, len(c.UnformattedUrls))
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()
	for i, unformattedUrl := range c.UnformattedUrls {
		name := c.Names[i]
		url := fmt.Sprintf(unformattedUrl, n)
		wg.Add(1)
		go c.getFromUrl(ctx, ch, wg, name, url)
	}

	wgCh := make(chan byte)
	go func() {
		wg.Wait()
		wgCh <- 0
	}()

	select {
	case r := <-ch:
		return r, nil
	case <-wgCh:
		return nil, errors.New("all requests failed")
	case <-time.Tick(time.Second):
		return nil, errors.New("timeout")
	}
}

func (c *Cep) getFromUrl(ctx context.Context, ch chan<- []byte, wg *sync.WaitGroup, name, url string) {
	defer wg.Done()
	init := time.Now()
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		log.Printf("invalid response from %s: %s", name, string(b))
		return
	}

	ch <- b
	log.Printf("valid response from %s(%dms): \n%s", name, time.Now().Sub(init).Milliseconds(), string(b))
}
