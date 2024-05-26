package main

import (
	"context"
	"fmt"
	"math/big"
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

func InitializeClient(config *Config) (*ethclient.Client, *bind.TransactOpts, error) {
	// Декодируем приватный ключ
	pk, err := crypto.HexToECDSA(config.PrivateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode private key: %v", err)
	}

	// Подключаемся к RPC клиенту
	client, err := ethclient.Dial(config.RPCURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the Ethereum client: %v", err)
	}

	// Получаем Chain ID
	var chainIDHex string
	rpcClient, err := rpc.Dial(config.RPCURL)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to connect to the Ethereum RPC client: %v", err)
	}
	err = rpcClient.Call(&chainIDHex, "eth_chainId")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get chain ID: %v", err)
	}

	// Преобразуем Chain ID из строки в *big.Int
	chainID, err := hexutil.DecodeBig(chainIDHex)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to decode chain ID: %v", err)
	}

	// Создаем авторизованный транзактор
	auth, err := bind.NewKeyedTransactorWithChainID(pk, chainID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create authorized transactor: %v", err)
	}

	return client, auth, nil
}

func SwapExactETHForTokens(client *ethclient.Client, auth *bind.TransactOpts, config *Config, amount int, slippage int) error {
	// Загружаем ABI контракта Uniswap V2 Router
	parsedABI, err := abi.JSON(strings.NewReader(config.UniswapV2RouterABI))
	if err != nil {
		return fmt.Errorf("failed to parse Uniswap V2 Router ABI: %v", err)
	}

	// Устанавливаем минимальное количество токенов, которое мы готовы получить (с учетом проскальзывания)
	amountOutMin := big.NewInt(int64(amount * (100 - slippage) / 100))

	// Устанавливаем путь свопа (WETH -> USDT)
	path := []common.Address{
		config.WETHAddress, // WETH
		config.USDTAddress,
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
	gasPrice := big.NewInt(16000000000) // 16 Gwei

	// Захардкодим лимит газа (gas limit)
	gasLimit := uint64(700000)

	// Создаем данные для вызова функции swapExactETHForTokens
	input, err := parsedABI.Pack("swapExactETHForTokens", amountOutMin, path, to, deadline)
	if err != nil {
		return fmt.Errorf("failed to pack input data: %v", err)
	}

	// Создаем транзакцию
	tx := types.NewTransaction(nonce, config.UniswapV2RouterAddress, big.NewInt(int64(amount)*big.NewInt(1e18).Int64()), gasLimit, gasPrice, input)

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
