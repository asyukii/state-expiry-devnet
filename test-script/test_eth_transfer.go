// This file is used to test ETH transfer. The targetted account will send ETH to some random account.
// Users can specify the TPS when running this program. It will keep sending ETH until it is interrupted.

package main

import (
	"context"
	"crypto/ecdsa"
	"flag"
	"fmt"
	"github.com/bnb-chain/bsc-deploy/test-script/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"math/big"
	"time"
)

var edpoint = "http://localhost:8503"
var chainId = big.NewInt(123454321)

func sendEther(client *ethclient.Client, key *ecdsa.PrivateKey, toAddr common.Address, value *big.Int, nonce uint64) (common.Hash, error) {
	gasLimit := uint64(3e4)
	gasPrice := big.NewInt(params.GWei * 10000)

	tx := types.NewTransaction(nonce, toAddr, value, gasLimit, gasPrice, nil)
	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainId), key)
	if err != nil {
		return common.Hash{}, fmt.Errorf("sign tx failed, %v", err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("send tx failed, %v", err)
	}
	txhash := signedTx.Hash()
	return txhash, nil
}

func main() {
	var tps uint

	// Prepare flag parameters
	flag.UintVar(&tps, "tps", 50, "Transaction-Per-Second (TPS)")
	flag.Parse()

	senderPrvKey := utils.ParsePrivateKey("190e410a96c56dcc7cbe6ee04ce68fbcf2eb7d86c441e840235373078cf6bb0c")
	senderAddr := crypto.PubkeyToAddress(senderPrvKey.PublicKey) // 0xD8C0Aa483406A1891E5e03B21F2bc01379fc3b20

	c, _ := ethclient.Dial(edpoint)
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
				nonce, err := c.PendingNonceAt(context.Background(), senderAddr)
				if err != nil {
					fmt.Println(err)
					continue
				}
				_, err = sendEther(c, senderPrvKey, receiverAddr, big.NewInt(params.GWei*1), nonce)
				if err != nil {
					fmt.Println(err)
					continue
				}
			}
			fmt.Printf("sent %d transactions\n", tps)
		}
	}
}
