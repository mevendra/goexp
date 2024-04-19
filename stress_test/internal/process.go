package internal

import (
	"log"
	"net/http"
)

func process(url string, ch chan<- result) {
	for {
		res, err := http.Get(url)
		if err != nil {
			log.Println("unexpected error processing request:", err)
			ch <- result{Err: err}
			continue
		}
		if res.Body != nil {
			_ = res.Body.Close()
		}
		ch <- result{StatusCode: res.StatusCode}
	}
}
