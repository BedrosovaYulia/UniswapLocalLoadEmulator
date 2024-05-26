package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"
)

func main() {
	config, err := LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	client, auth, err := InitializeClient(config)
	if err != nil {
		log.Fatalf("Failed to initialize client: %v", err)
	}

	fmt.Println("Authorized transactor created successfully")

	for {
		amount := rand.Intn(10) + 1
		slippage := rand.Intn(46) + 5

		fmt.Printf("Buying %d USDT with ETH with %d%% slippage\n", amount, slippage)

		// Выполняем своп
		err := SwapExactETHForTokens(client, auth, config, amount, slippage)
		if err != nil {
			log.Printf("Failed to perform swap: %v", err)
		}

		time.Sleep(5 * time.Second)
	}
}
