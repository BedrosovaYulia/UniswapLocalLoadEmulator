package main

import (
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

type Config struct {
	RPCURL                 string
	PrivateKey             string
	UniswapV2RouterABI     string
	UniswapV2RouterAddress common.Address
	USDTAddress            common.Address
	WETHAddress            common.Address
}

func LoadConfig() (*Config, error) {
	rpcURL := os.Getenv("RPC_URL")
	privateKey := os.Getenv("PRIVATE_KEY")

	if rpcURL == "" || privateKey == "" {
		log.Fatal("Error loading environment variables")
	}

	// Удаляем префикс "0x" из приватного ключа, если он есть
	privateKey = strings.TrimPrefix(privateKey, "0x")

	return &Config{
		RPCURL:                 rpcURL,
		PrivateKey:             privateKey,
		UniswapV2RouterABI:     `[{"constant":false,"inputs":[{"internalType":"uint256","name":"amountOutMin","type":"uint256"},{"internalType":"address[]","name":"path","type":"address[]"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"swapExactETHForTokens","outputs":[{"internalType":"uint256[]","name":"amounts","type":"uint256[]"}],"payable":true,"stateMutability":"payable","type":"function"}]`,
		UniswapV2RouterAddress: common.HexToAddress("0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"),
		USDTAddress:            common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"),
		WETHAddress:            common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"),
	}, nil
}
