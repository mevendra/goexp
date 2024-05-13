package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	log.Println("Starting client")

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	err = db.AutoMigrate(&QuotationDTO{})
	if err != nil {
		panic(err)
	}

	quotation := NewQuotation(time.Millisecond*200, time.Millisecond*10, db)
	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		res, err := quotation.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		b, err := json.Marshal(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

type QuotationDatabase struct {
	timeout time.Duration
	db      *gorm.DB
}

func NewQuotationDatabase(timeout time.Duration, db *gorm.DB) QuotationDatabase {
	return QuotationDatabase{timeout: timeout, db: db}
}

func (q QuotationDatabase) Add(quotation QuotationDTO) {
	ctx, fnCancel := context.WithTimeout(context.Background(), q.timeout)
	defer fnCancel()

	aux := q.db.WithContext(ctx)
	err := aux.Create(&quotation).Error
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("database timeout exceeded")
		} else {
			log.Println("Error executing statement", err)
		}
		return
	}
}

type Quotation struct {
	timeoutHttp time.Duration
	quotationDB QuotationDatabase
}

func NewQuotation(timeoutHttp, timeoutDb time.Duration, db *gorm.DB) Quotation {
	quotationDb := NewQuotationDatabase(timeoutDb, db)
	return Quotation{
		timeoutHttp: timeoutHttp,
		quotationDB: quotationDb,
	}
}

func (q Quotation) New() (QuotationDTO, error) {
	ctx, fnCancel := context.WithTimeout(context.Background(), q.timeoutHttp)
	defer fnCancel()
	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		return QuotationDTO{}, err
	}
	client := &http.Client{}

	res, err := client.Do(req)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			log.Println("request timeout exceeded")
		}
		return QuotationDTO{}, err
	}

	if res.StatusCode != 200 {
		return QuotationDTO{}, fmt.Errorf("invalid status code: %d", res.StatusCode)
	}

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return QuotationDTO{}, err
	}

	var quotationResponse struct {
		UsdBrl struct {
			Bid string `json:"bid"`
		} `json:"USDBRL"`
	}
	err = json.Unmarshal(b, &quotationResponse)
	if err != nil {
		return QuotationDTO{}, err
	}

	bid, err := strconv.ParseFloat(quotationResponse.UsdBrl.Bid, 64)
	if err != nil {
		return QuotationDTO{}, err
	}
	dto := QuotationDTO{
		Bid: bid,
	}
	q.quotationDB.Add(dto)
	return dto, nil
}

type QuotationDTO struct {
	Id        int       `json:"id"`
	Bid       float64   `json:"bid"`
	CreatedAt time.Time `json:"created_at"`
}
