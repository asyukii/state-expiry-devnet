// This file is used to test ERC20 token transfer with expiry. The targetted account will send tokens to some random account.
// Users can specify the TPS and expiry percentage when running this program.
// It will keep sending tokens to unexpired accounts and periodically send to expired accounts until it is interrupted.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bnb-chain/bsc-deploy/test-script/utils"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"crypto/ecdsa"
	"log"
	"math/big"
	"time"
)

var accountNum = 3_000

func createBatchAccounts(expirePerc, tps uint) []common.Address {

	accounts := make([]common.Address, accountNum)

	for i := 0; i < accountNum; i++ {
		prvKey, err := crypto.GenerateKey()
		if err != nil {
			fmt.Println("got err when GenerateKey", err)
			continue
		}
		accounts[i] = crypto.PubkeyToAddress(prvKey.PublicKey)
	}

	return accounts
}

func sendTransactionsToAccounts(client *ethclient.Client, senderPrvKey *ecdsa.PrivateKey, nonce uint64, gasPrice *big.Int, contract common.Address, abi abi.ABI, accounts []common.Address) (uint64, error) {

	for _, receiverAddr := range accounts {
		num, _ := new(big.Int).SetString("1000000000000000000", 10)
		input, err := abi.Pack("transfer", receiverAddr, num)
		if err != nil {
			fmt.Println("got err when Pack", err)
			continue
		}

		// SendTransaction
		tx := types.NewTx(&types.LegacyTx{
			Nonce:    nonce,
			GasPrice: gasPrice,
			Gas:      uint64(200000),
			To:       &contract,
			Data:     input,
		})

		chainId := big.NewInt(123454321)
		signedTx, err := types.SignTx(tx, types.NewLondonSigner(chainId), senderPrvKey)
		if err != nil {
			fmt.Println("got err when SignTx", err)
			continue
		}

		_, err = types.Sender(types.NewLondonSigner(chainId), signedTx)
		if err != nil {
			fmt.Println("got err when Sender", err)
			continue
		}

		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		err = client.SendTransaction(ctx, signedTx)
		if err != nil {
			fmt.Println("got err when SendTransaction", err)
			continue
		}

		fmt.Printf("txHash %v\n", signedTx.Hash())
		nonce++
	}

	fmt.Printf("sent %d transactions\n", len(accounts))

	return nonce, nil
}

func main() {

	var tps uint
	var expirePerc uint

	// Prepare flag parameters
	flag.UintVar(&tps, "tps", 50, "Transaction-Per-Second (TPS)")
	flag.UintVar(&expirePerc, "expirePerc", 5, "Percentage of expired slot")
	flag.Parse()

	// Get contract address and abi
	contracts := utils.ReadDeployedContracts("../test-contract/deployed_contracts.json")
	contract, ok := contracts["ERC20Token"]
	if !ok {
		log.Fatal("cannot find ERC20Token contract address")
	}
	erc20 := utils.LoadAbi("abi/ERC20Token.json")

	// Get sender address
	senderPrvKey := utils.ParsePrivateKey("190e410a96c56dcc7cbe6ee04ce68fbcf2eb7d86c441e840235373078cf6bb0c")
	senderAddr := crypto.PubkeyToAddress(senderPrvKey.PublicKey)

	// Initialize client
	client, err := ethclient.Dial("http://localhost:8503")
	utils.Fatal(err)
	defer client.Close()

	// Get nonce and gas price
	gasPrice, err := client.SuggestGasPrice(context.Background())
	utils.Fatal(err)
	nonce, err := client.PendingNonceAt(context.Background(), senderAddr)
	utils.Fatal(err)

	// Create accounts and send transactions to them
	accounts := createBatchAccounts(expirePerc, tps)
	nonce, err = sendTransactionsToAccounts(client, senderPrvKey, nonce, gasPrice, contract, erc20, accounts) // for the first time, we send transactions to all accounts
	if err != nil {
		fmt.Println("got err when sendTransactionsToAccounts", err)
		return
	}
	time.Sleep(180 * time.Second) // wait for 3 minutes so that accounts get expired

	expiredCount := tps * expirePerc / 100
	unexpiredCount := tps - expiredCount
	nonce, err = sendTransactionsToAccounts(client, senderPrvKey, nonce, gasPrice, contract, erc20, accounts[:unexpiredCount]) // for the second time, we send transactions to unexpired accounts
	if err != nil {
		fmt.Println("got err when sendTransactionsToAccounts", err)
		return
	}
	time.Sleep(1000 * time.Millisecond)

	startIndex := unexpiredCount
	endIndex := startIndex + expiredCount
	t := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-t.C:

			// If we have sent transactions to all accounts, we need to reset the index
			// If we start from the beginning, the accounts should be expired again
			if int(endIndex) > len(accounts)-int(unexpiredCount) {
				startIndex = unexpiredCount
				endIndex = startIndex + expiredCount
			}

			// Always send transactions to unexpired accounts
			nonce, err = sendTransactionsToAccounts(client, senderPrvKey, nonce, gasPrice, contract, erc20, accounts[:unexpiredCount])
			if err != nil {
				fmt.Println("got err when sendTransactionsToAccounts", err)
				return
			}

			// Send transactions to expired accounts periodically
			nonce, err = sendTransactionsToAccounts(client, senderPrvKey, nonce, gasPrice, contract, erc20, accounts[startIndex:endIndex])
			startIndex = endIndex
			endIndex = startIndex + expiredCount
		}
	}
}
