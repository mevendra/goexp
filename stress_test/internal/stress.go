package internal

import (
	"fmt"
	"log"
	netUrl "net/url"
	"time"
)

func Run(url string, requests, concurrency int) {
	// Check for correct values
	if concurrency <= 0 {
		panic("concurrency must be > 0")
	}
	if requests <= 0 {
		panic("requests must be > 0")
	}
	_, err := netUrl.ParseRequestURI(url)
	if err != nil {
		panic(fmt.Sprintf("URL parse failed: %s", err))
	}

	log.Printf("Initializing Stress Test with [URL: %s, Requests: %d, Concurrency: %d]\n", url, requests, concurrency)
	ch := make(chan result)
	for i := 0; i < concurrency; i++ {
		go process(url, ch)
	}
	start := time.Now()
	report := Report{Others: make(map[int]int)}
	for res := range ch {
		if res.Err != nil {
			report.numberOfInternalErrors++
			if report.numberOfInternalErrors >= requests {
				panic(fmt.Sprintf("%d unexpected errors, last error: %s", report.numberOfInternalErrors, res.Err))
			}
			continue
		}

		report.NumberOfRequests++
		if res.StatusCode == 200 {
			report.NumberOfSuccess++
		} else {
			report.Others[res.StatusCode]++
		}
		if report.NumberOfRequests >= requests {
			break
		}
	}
	report.TotalTime = time.Since(start).Milliseconds()
	log.Println(report.String())
}
