package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"math/big"
	"os"

	uniETHToken "simple-token-wallet/token"
	"simple-token-wallet/util"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/urfave/cli/v2"
)

func printBalance(balance *big.Int, balanceType string) {
	fbalance := new(big.Float)
	fbalance.SetString(balance.String())
	ethBalance := new(big.Float).Quo(fbalance, big.NewFloat(math.Pow10(18)))
	fmt.Printf("%v Balance: %v wei (%v %v)\n", balanceType, balance, ethBalance, balanceType)
}

func getBalance(client *ethclient.Client, addr string) *big.Int {
	accountAddress := common.HexToAddress(addr)
	balance, err := client.BalanceAt(context.Background(), accountAddress, nil)
	if err != nil {
		log.Fatalf("Error getting balance: %v", err)
	}

	return balance
}

func getUniETHBalance(client *ethclient.Client, tokenAddr string, addr string) (balance *big.Int) {
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

	app := &cli.App{
		Name:  "Simple Token Wallet",
		Usage: "A simple CLI wallet for balance and sending tokens",
		Commands: []*cli.Command{
			{
				Name:  "balance",
				Usage: "Check balance for a given address",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "address",
						Aliases:  []string{"a"},
						Usage:    "The wallet `ADDRESS` to check the balance of",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "apiKey",
						Aliases:  []string{"z"},
						Usage:    "Access Node `ApiKey` can be retrieved from https://access.rockx.com",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					address := c.String("address")
					apiKey := c.String("apiKey")

					client, err := ethclient.Dial(fmt.Sprintf("https://eth.w3node.com/%v/api", apiKey))

					if err != nil {
						log.Fatalf("Failed to connect to eth client: %v", err)
					}
					printBalance(getBalance(client, address), "ETH")
					return nil
				},
			},
			{
				Name:  "uniethBalance",
				Usage: "Checks the uniETH balance for the given address",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "address",
						Aliases:  []string{"a"},
						Usage:    "The wallet `Address` to check the balance of",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "apikey",
						Aliases:  []string{"z"},
						Usage:    "Access Node `ApiKey` can be retrieved from https://access.rockx.com",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					address := c.String("address")
					apiKey := c.String("apikey")

					config, err := util.LoadConfig(".")
					if err != nil {
						log.Fatalf("Failed to load config: %v", err)
					}

					client, err := ethclient.Dial(fmt.Sprintf("https://eth.w3node.com/%v/api", apiKey))
					if err != nil {
						log.Fatalf("Failed to connect to eth client: %v", err)
					}
					printBalance(getUniETHBalance(client, config.UniethTokenAddress, address), "uniETH")
					return nil
				},
			},
			{
				Name:  "chainID",
				Usage: "Returns the chain ID",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "apikey",
						Aliases:  []string{"z"},
						Usage:    "Access Node `ApiKey` can be retrieved from https://access.rockx.com",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					apiKey := c.String("apikey")

					client, err := ethclient.Dial(fmt.Sprintf("https://eth.w3node.com/%v/api", apiKey))
					if err != nil {
						log.Fatalf("Failed to connect to eth client: %v", err)
					}

					chainID := getChainID(client)
					fmt.Printf("Chain ID: %v\n", chainID)

					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
