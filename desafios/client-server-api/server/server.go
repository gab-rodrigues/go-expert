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
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Exchange struct {
	Code       string `json:"code"`
	Codein     string `json:"codein"`
	Name       string `json:"name"`
	High       string `json:"high"`
	Low        string `json:"low"`
	VarBid     string `json:"varBid"`
	PctChange  string `json:"pctChange"`
	Bid        string `json:"bid"`
	Ask        string `json:"ask"`
	Timestamp  string `json:"timestamp"`
	CreateDate string `json:"create_date"`
}

type Response struct {
	USDBRL Exchange `json:"USDBRL"`
}

func main() {
	http.HandleFunc("/cotacao", Handle)
	http.ListenAndServe(":8080", nil)
}

func Handle(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
		http.Error(w, "Cancelled request by client", http.StatusRequestTimeout)
		return
	default:
	}

	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", "https://economia.awesomeapi.com.br/json/last/USD-BRL", nil)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("Request timed out to external client")
		}
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	exchangeJson, err := io.ReadAll(resp.Body)
	exchange, err := saveExchange(exchangeJson)
	if err != nil {
		http.Error(w, "Internal server error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	println("response do endpoint")
	println(string(exchangeJson))

	response, _ := json.Marshal(exchange.Bid)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}

func saveExchange(exchangeJson []byte) (*Exchange, error) {
	os.Mkdir("./data", os.ModePerm)

	db, err := gorm.Open(sqlite.Open("./data/data.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect database:", err)
		return nil, err
	}

	db.AutoMigrate(&Exchange{})

	var respHttpPost Response
	if err := json.Unmarshal(exchangeJson, &respHttpPost); err != nil {
		panic(err)
		return nil, err
	}

	fmt.Println(respHttpPost)
	fmt.Printf("Bid apos unmarshall %s \n", respHttpPost.USDBRL.Bid)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()

	if err := db.WithContext(ctx).Create(&respHttpPost.USDBRL).Error; err != nil {
		if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			log.Println("database operation timed out")
			return nil, fmt.Errorf("save operation timed out")
		}
		return nil, fmt.Errorf("failed to insert data: %w", err)
	}

	return &respHttpPost.USDBRL, nil
}
