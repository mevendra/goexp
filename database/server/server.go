package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func main() {
	log.Println("Starting client")

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	_, err = db.Query("CREATE TABLE IF NOT EXISTS bids (" +
		"id INT PRIMARY KEY AUTO_INCREMENT," +
		"bid FLOAT NOT NULL," +
		"created_at TIME NOT NULL" +
		");")
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
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write(b)
	})
	_ = http.ListenAndServe(":8080", nil)
}

type QuotationDatabase struct {
	timeout time.Duration
	db      *sql.DB
}

func NewQuotationDatabase(timeout time.Duration, db *sql.DB) QuotationDatabase {
	return QuotationDatabase{timeout: timeout, db: db}
}

func (q QuotationDatabase) Add(quotation QuotationDTO) {
	stmt, err := q.db.Prepare("insert into bids(created_at, bid) values(?, ?)")
	if err != nil {
		log.Println("Error preparing statement", err)
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now(), quotation.Bid)
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

func NewQuotation(timeoutHttp, timeoutDb time.Duration, db *sql.DB) Quotation {
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
	Bid float64 `json:"bid"`
}
