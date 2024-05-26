package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	rpcURL := os.Getenv("RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")

	if rpcURL == "" || privateKey == "" {
		log.Fatal("Error loading environment variables")
	}

	fmt.Printf("RPC_URL: %s\n", rpcURL)
	fmt.Printf("PRIVATE_KEY: %s\n", privateKey)

	// Удаляем префикс "0x" из приватного ключа, если он есть
	if len(privateKey) > 2 && privateKey[:2] == "0x" {
		privateKey = privateKey[2:]
	}

	// Декодируем приватный ключ
	pk, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatalf("Failed to decode private key: %v", err)
	}

	// Подключаемся к RPC клиенту
	client, err := rpc.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Получаем Chain ID
	var chainIDHex string
	err = client.Call(&chainIDHex, "eth_chainId")
	if err != nil {
		log.Fatalf("Failed to get chain ID: %v", err)
	}

	// Преобразуем Chain ID из строки в *big.Int
	chainID, err := hexutil.DecodeBig(chainIDHex)
	if err != nil {
		log.Fatalf("Failed to decode chain ID: %v", err)
	}

	// Создаем авторизованный транзактор
	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}

	_ = auth

	fmt.Println("Authorized transactor created successfully")

	for {
		amount := rand.Intn(10) + 1
		slippage := rand.Intn(46) + 5

		fmt.Printf("Buying %d USDT with ETH with %d%% slippage\n", amount, slippage)

		// Here you would add the logic to interact with Uniswap V2 contracts
		// This is a placeholder for the actual swap logic
		// Example: uniswapV2Router.SwapExactETHForTokens(...)

		time.Sleep(5 * time.Second)
	}
}
