package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/gorm"
)

type CurrencyConversionResponse struct {
	gorm.Model `json:"-"`
	Value      float64 `json:"value"`
}

const (
	clientTimout  time.Duration = 300 * time.Millisecond
	serverAddress string        = "http://localhost:8080/cotacao"
)

func writeDollarValueToTxT(value float64) error {

	f, err := os.Create("cotacao.txt")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(fmt.Sprintf("Dólar: R$ %.2f\n", value))
	if err != nil {
		return err
	}
	return nil
}

func RequestCurrencyConversion() error {

	var serverResponse CurrencyConversionResponse

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, clientTimout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", serverAddress, nil)
	if err != nil {
		return err
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Request timed out")
		}
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("error: %s", res.Status)
	}

	err = json.NewDecoder(res.Body).Decode(&serverResponse)
	if err != nil {
		return err
	}
	log.Printf("Valor do Dólar: R$ %.2f\n", serverResponse.Value)
	err = writeDollarValueToTxT(serverResponse.Value)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	err := RequestCurrencyConversion()
	if err != nil {
		log.Fatalf("Erro ao fazer a requisição: %v", err)
	} else {
		log.Println("Requisição realizada com sucesso ")
	}
}
