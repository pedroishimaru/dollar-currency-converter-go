package main

import (
	"desafio"
	"log"
	"time"
)

func main() {

	go desafio.StartServer()
	time.Sleep(time.Second)

	_ = desafio.RequestCurrencyConversion()

	err := desafio.RequestCurrencyConversion()
	if err != nil {
		panic(err)
	} else {
		log.Println("Requisição realizada com sucesso ")
	}

}
