package desafio

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"time"
)

type currencyConvertion struct {
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

const (
	apiURL       string        = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	dbTimeout    time.Duration = 10 * time.Millisecond
	queryTimeout time.Duration = 200 * time.Millisecond
)

func GetValue(ctx context.Context) (float64, error) {

	valueCtx, cancel := context.WithTimeout(ctx, queryTimeout)
	defer cancel()

	var currencyData currencyConvertion
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
		http.Error(w, "Error getting value", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]float64{"value": value}
	json.NewEncoder(w).Encode(response)

}

func StartServer() {

	http.HandleFunc("/cotacao", handler)
	http.ListenAndServe(":8080", nil)

}
