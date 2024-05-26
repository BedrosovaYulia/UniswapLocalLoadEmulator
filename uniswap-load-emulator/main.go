package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	uniswapV2RouterABI     = `[{"constant":false,"inputs":[{"internalType":"uint256","name":"amountOutMin","type":"uint256"},{"internalType":"address[]","name":"path","type":"address[]"},{"internalType":"address","name":"to","type":"address"},{"internalType":"uint256","name":"deadline","type":"uint256"}],"name":"swapExactETHForTokens","outputs":[{"internalType":"uint256[]","name":"amounts","type":"uint256[]"}],"payable":true,"stateMutability":"payable","type":"function"}]`
	uniswapV2RouterAddress = "0x7a250d5630B4cF539739dF2C5dAcb4c659F2488D"
	usdtAddress            = "0xdAC17F958D2ee523a2206206994597C13D831ec7" // USDT contract address on Ethereum mainnet
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
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Получаем Chain ID
	var chainIDHex string
	rpcClient, err := rpc.Dial(rpcURL)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum RPC client: %v", err)
	}
	err = rpcClient.Call(&chainIDHex, "eth_chainId")
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

	// Загружаем ABI контракта Uniswap V2 Router
	parsedABI, err := abi.JSON(strings.NewReader(uniswapV2RouterABI))
	if err != nil {
		log.Fatalf("Failed to parse Uniswap V2 Router ABI: %v", err)
	}

	routerAddress := common.HexToAddress(uniswapV2RouterAddress)
	usdt := common.HexToAddress(usdtAddress)

	fmt.Println("Authorized transactor created successfully")

	for {
		amount := rand.Intn(10) + 1
		slippage := rand.Intn(46) + 5

		fmt.Printf("Buying %d USDT with ETH with %d%% slippage\n", amount, slippage)

		// Выполняем своп
		err := swapExactETHForTokens(client, auth, parsedABI, routerAddress, usdt, amount, slippage)
		if err != nil {
			log.Printf("Failed to perform swap: %v", err)
		}

		time.Sleep(5 * time.Second)
	}
}

func swapExactETHForTokens(client *ethclient.Client, auth *bind.TransactOpts, parsedABI abi.ABI, routerAddress common.Address, tokenAddress common.Address, amount int, slippage int) error {
	// Устанавливаем минимальное количество токенов, которое мы готовы получить (с учетом проскальзывания)
	amountOutMin := big.NewInt(int64(amount * (100 - slippage) / 100))

	// Устанавливаем путь свопа (ETH -> USDT)
	path := []common.Address{
		common.HexToAddress("0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"), // ETH
		tokenAddress,
	}

	// Устанавливаем адрес получателя и дедлайн
	to := auth.From
	deadline := big.NewInt(time.Now().Add(15 * time.Minute).Unix())

	// Получаем текущий nonce для аккаунта
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		return fmt.Errorf("failed to get account nonce: %v", err)
	}

	// Захардкодим цену газа (gas price)
	gasPrice := big.NewInt(16000000000) // 8 Gwei

	// Захардкодим лимит газа (gas limit)
	gasLimit := uint64(60000)

	// Создаем данные для вызова функции swapExactETHForTokens
	input, err := parsedABI.Pack("swapExactETHForTokens", amountOutMin, path, to, deadline)
	if err != nil {
		return fmt.Errorf("failed to pack input data: %v", err)
	}

	// Создаем транзакцию
	tx := types.NewTransaction(nonce, routerAddress, big.NewInt(int64(amount)*big.NewInt(1e18).Int64()), gasLimit, gasPrice, input)

	// Подписываем транзакцию
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Отправляем транзакцию
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return fmt.Errorf("failed to send transaction: %v", err)
	}

	fmt.Printf("Transaction sent: %s\n", signedTx.Hash().Hex())
	return nil
}
