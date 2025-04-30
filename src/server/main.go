package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type currencyConversion struct {
	USDBRL struct {
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
	} `json:"USDBRL"`
}

type CurrencyConversionResponse struct {
	gorm.Model `json:"-"`
	Value      float64 `json:"value"`
}

const (
	apiURL       string        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbTimeout    time.Duration = 10 * time.Millisecond
	queryTimeout time.Duration = 200 * time.Millisecond
)

var db *gorm.DB

func SetupDatabase() (*gorm.DB, error) {
	// Open SQLite database (it will create the file if it doesn't exist)
	db, err := gorm.Open(sqlite.Open("conversions.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// Auto migrate the schema
	err = db.AutoMigrate(&CurrencyConversionResponse{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func InsertCurrencyConversion(baseCtx context.Context, conversion CurrencyConversionResponse) error {

	ctx, cancel := context.WithTimeout(baseCtx, dbTimeout)
	defer cancel()

	result := db.WithContext(ctx).Create(&conversion)

	if result.Error != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Database operation timed out")
		}
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil

}

func GetValue(ctx context.Context) (float64, error) {

	valueCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var currencyData currencyConversion
	c := http.Client{}

	req, err := http.NewRequestWithContext(valueCtx, http.MethodGet, apiURL, nil)
	if err != nil {
		return -1, err
	}

	resp, err := c.Do(req)
	if err != nil {
		return -1, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, err
	}

	err = json.Unmarshal(body, &currencyData)
	if err != nil {
		return -1, err
	}

	value, err := strconv.ParseFloat(currencyData.USDBRL.Bid, 64)
	if err != nil {
		return -1, err
	}

	return value, nil

}

func handler(w http.ResponseWriter, r *http.Request) {

	baseCtx := r.Context()

	value, err := GetValue(baseCtx)
	if err != nil {

		log.Println("Error getting value:", err)

		if err == context.DeadlineExceeded {
			log.Println("Request timed out")
			http.Error(w, "Request timed out", http.StatusGatewayTimeout)
			return
		} else {
			http.Error(w, "Error getting value", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := CurrencyConversionResponse{Value: value}
	err = InsertCurrencyConversion(baseCtx, response)
	if err != nil {
		log.Println("Error inserting conversion:", err)
		http.Error(w, "Error inserting conversion", http.StatusInternalServerError)
	} else {
		json.NewEncoder(w).Encode(response)
	}

}

func StartServer() {

	var err error

	log.Println("Starting server...")

	db, err = SetupDatabase()
	if err != nil {
		panic("Failed to connect to database")
	} else {
		log.Println("Connected to database")
	}

	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)

}

func main() {
	StartServer()
}
