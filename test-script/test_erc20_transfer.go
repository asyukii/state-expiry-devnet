// This file is used to test ERC20 token transfer. The targetted account will send tokens to some random account.
// Users can specify the TPS when running this program. It will keep sending tokens until it is interrupted.

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/bnb-chain/bsc-deploy/test-script/utils"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func main() {

	var tps uint

	// Prepare flag parameters
	flag.UintVar(&tps, "tps", 50, "Transaction-Per-Second (TPS)")
	flag.Parse()

	contracts := utils.ReadDeployedContracts("../test-contract/deployed_contracts.json")
	contract, ok := contracts["ERC20Token"]
	if !ok {
		log.Fatal("cannot find ERC20Token contract address")
	}
	senderPrvKey := utils.ParsePrivateKey("190e410a96c56dcc7cbe6ee04ce68fbcf2eb7d86c441e840235373078cf6bb0c")
	senderAddr := crypto.PubkeyToAddress(senderPrvKey.PublicKey)

	erc20 := utils.LoadAbi("abi/ERC20Token.json")

	client, err := ethclient.Dial("http://localhost:8503")
	utils.Fatal(err)
	defer client.Close()

	gasPrice, err := client.SuggestGasPrice(context.Background())
	utils.Fatal(err)
	nonce, err := client.PendingNonceAt(context.Background(), senderAddr)
	utils.Fatal(err)

	t := time.NewTicker(1000 * time.Millisecond)
	for {
		select {
		case <-t.C:
			for i := uint(0); i < tps; i++ {
				prvKey, err := crypto.GenerateKey()
				if err != nil {
					fmt.Println("got err when GenerateKey", err)
					continue
				}
				receiverAddr := crypto.PubkeyToAddress(prvKey.PublicKey)
				fmt.Printf("newAccount %v prvKey %v\n", receiverAddr, prvKey.D.String())

				num, _ := new(big.Int).SetString("1000000000000000000", 10)
				input, err := erc20.Pack("transfer", receiverAddr, num)
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
			fmt.Printf("sent %d transactions\n", tps)
		}
	}
}
