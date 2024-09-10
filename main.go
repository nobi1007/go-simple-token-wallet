package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	uniETHToken "simple-token-wallet/token"
	"simple-token-wallet/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func getBalance(client *ethclient.Client, addr string) (*big.Int, *big.Float) {
	accountAddress := common.HexToAddress(addr)
	balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		log.Fatalf("Error getting balance: %v", err)
	}

	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18))).SetPrec(8)

	return balance, ethValue
}

func getUniETHBalance(client *ethclient.Client, tokenAddr string, addr string) (balance *big.Int, ethBalance *big.Float) {
	tokenAddress := common.HexToAddress(tokenAddr)
	accountAddress := common.HexToAddress(addr)

	tokenInstance, err := uniETHToken.NewToken(tokenAddress, client)
	if err != nil {
		log.Fatalf("Failed to load the token contract: %v", err)
	}

	balance, err = tokenInstance.BalanceOf(&bind.CallOpts{}, accountAddress)
	if err != nil {
		log.Fatalf("Failed to load the uniETH balance: %v\n", err)
	}
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethBalance = new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	return
}

func getChainID(client *ethclient.Client) *big.Int {
	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Error getting chain ID: %v\n", err)
	}
	return chainID
}

func main() {

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatalf("Failed to load config: %v\n", err)
	}

	client, err := ethclient.Dial(fmt.Sprintf("https://eth.w3node.com/%v/api", config.EthClientAPIKey))
	if err != nil {
		log.Fatalf("Failed to connect to eth client: %v", err)
	}
	fmt.Println("Successfully connected to eth client")

	fmt.Printf("Chain ID: %v\n", getChainID(client))

	weiBalance, ethBalance := getBalance(client, config.UserAddress)

	fmt.Printf("ETH Balance: %v wei (%v ETH)\n", weiBalance, ethBalance)

	uniEthWeiBalance, uniEthBalance := getUniETHBalance(client, config.UniethTokenAddress, config.UserAddress)

	fmt.Printf("uniETH Balance: %v wei (%v uniETH)\n", uniEthWeiBalance, uniEthBalance)
}
