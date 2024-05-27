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

		// Получаем баланс пользователя в ETH до свопа
		ethBalanceBefore, err := GetETHBalance(client, auth.From)
		if err != nil {
			log.Printf("Failed to get ETH balance: %v", err)
			continue
		}

		fmt.Printf("Balance before swap: ETH = %s\n", ethBalanceBefore.String())

		// Выполняем своп
		err = SwapExactETHForTokens(client, auth, config, amount, slippage)
		if err != nil {
			log.Printf("Failed to perform swap: %v", err)
			continue
		}

		// Получаем баланс пользователя в ETH после свопа
		ethBalanceAfter, err := GetETHBalance(client, auth.From)
		if err != nil {
			log.Printf("Failed to get ETH balance: %v", err)
			continue
		}

		fmt.Printf("Balance after swap: ETH = %s\n", ethBalanceAfter.String())

		time.Sleep(5 * time.Second)
	}
}
