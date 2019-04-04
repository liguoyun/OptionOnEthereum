package main

import (
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

const (
	remote_api_address   = "https://ropsten.infura.io/v3/9cef3fa3f49847ec9a891cfa25ec58db"
	eth_contract_address = "0x3709e48417e6b59cc40d5d75606a81be949dd2f7"
	GasPrice             = 1000000000
	GasLimit             = 60000
	Decimal              = 1000
)

type StockTokenWapper struct {
	ApiAddr      string
	ContractAddr string
}

func NewStockTokenWapper(aAddr string, cAddr string) StockTokenWapper {
	stw := StockTokenWapper{
		ApiAddr:      aAddr,
		ContractAddr: cAddr,
	}
	return stw
}

func (s *StockTokenWapper) GetBalanceOfAddr(addr string) (string, error) {
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(s.ApiAddr)
	//fmt.Println("ethclient.Dial")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return "", err
	}
	// Instantiate the contract and display its name
	token, err := NewStockToken(common.HexToAddress(s.ContractAddr), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
		return "", err
	}
	banl, err := token.BalanceOf(nil, common.HexToAddress(addr))
	if err != nil {
		log.Fatalf("Failed to retrieve Balance: %v", err)
		return "", nil
	}
	bigAmount := big.NewInt(1)
	bigAmount = bigAmount.Div(banl, big.NewInt(Decimal))
	//fmt.Println("Balance is :", bigAmount)
	return bigAmount.String(), nil
}

func (s *StockTokenWapper) GetBalanceOfAddrs(userAddrs []string) []int64 {
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(s.ApiAddr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}
	// Instantiate the contract and display its name
	token, err := NewStockToken(common.HexToAddress(s.ContractAddr), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
	}
	var result []int64
	for _, addr := range userAddrs {
		banl, err := token.BalanceOf(nil, common.HexToAddress(addr))
		if err != nil {
			log.Fatalf("Failed to retrieve Balance: %v", err)
		}
		bigAmount := big.NewInt(1)
		bigAmount = bigAmount.Div(banl, big.NewInt(Decimal))
		fmt.Println("Balance is :", bigAmount)
		result = append(result, bigAmount.Int64())
	}
	return result
}

func GetStockTokenBalances(userAddrs []string) []int64 {
	contractor := NewStockTokenWapper(remote_api_address, eth_contract_address)
	return contractor.GetBalanceOfAddrs(userAddrs)
}

func (s *StockTokenWapper) TransferToken(fromAddr string, toAddr string, amount int64) (string, error) {
	// Create an IPC based RPC connection to a remote node
	conn, err := ethclient.Dial(s.ApiAddr)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
		return "", err
	}
	// Instantiate the contract and display its name
	token, err := NewStockToken(common.HexToAddress(s.ContractAddr), conn)
	if err != nil {
		log.Fatalf("Failed to instantiate a Token contract: %v", err)
		return "", err
	}

	//priv := "983F87D048AB7B50E64B3A06A5097591FBB8300F900297481709D288BF05601B"
	key, err := crypto.HexToECDSA(fromAddr)
	if err != nil {
		log.Fatalf("Failed to convert  private key: %v", err)
		return "", err
	}
	// Create an authorized transactor and spend 1 unicorn
	auth := bind.NewKeyedTransactor(key)
	auth.GasPrice = big.NewInt(GasPrice)
	auth.GasLimit = GasLimit
	bigAmount := big.NewInt(1)
	bigAmount = bigAmount.Mul(big.NewInt(amount), big.NewInt(Decimal))

	tx, err := token.Transfer(auth, common.HexToAddress(toAddr), bigAmount)
	if err != nil {
		log.Fatalf("Failed to request token transfer: %v", err)
		return "", err
	}
	fmt.Println("Transfer pending: ", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}
