package desafio

import (
	"os"
	"testing"
	"time"
)

func TestWriteTxt(t *testing.T) {
	err := writeDollarValueToTxT(5.26)
	if err != nil {
		t.Errorf("Failed to write to file: %v", err)
	}

	// Check if the file exists
	if _, err := os.Stat("cotacao.txt"); os.IsNotExist(err) {
		t.Errorf("File does not exist: %v", err)
	}
}

func TestClient(t *testing.T) {

	go StartServer()

	time.Sleep(1 * time.Second) // Give the server some time to start

	_ = RequestCurrencyConversion()

	err := RequestCurrencyConversion()
	if err != nil {
		t.Errorf("Failed to request currency conversion: %v", err)
	}
	// Check if the file exists
	if _, err := os.Stat("cotacao.txt"); os.IsNotExist(err) {
		t.Errorf("File does not exist: %v", err)
	}
}
