# State Expiry Devnet Setup Guide

## ⚙️ Installation

### Prerequisites

- Golang 1.19 or higher
- nodeJS 18 or higher

### Environment Setup

1. Clone repo and submodules.

```bash
git clone --recursive https://github.com/node-real/state-expiry-devnet.git
```

2. Prepare state expiry binary, and clone Geth repo with `state-expiry-mvp0.1` branch.

```bash
git clone https://github.com/0xbundler/go-ethereum --branch state-expiry-mvp0.1 geth-state-expiry
```

3. Compile `geth` and `bootnode`, then copy to `bin` folder

```bash
cd geth-state-expiry
make geth && go build -o bootnode ./cmd/bootnode
mkdir -p ../state-expiry-devnet/bin && cp build/bin/geth bootnode ../state-expiry-devnet/bin
```

## 🤖 Testing Flow

The complete flow of testing the state expiry features may look something like:

1. Deploy nodes
2. Deploy token contract
3. Read token balance
4. Wait for contract state to be expired
5. Read token balance again (expired error is returned)
6. Send revive transaction to revive contract state
7. Read token balance
8. Stop nodes

When a contract is deployed, you should wait for 2 epochs (100 blocks) so that its state is expired. Here's the table reference:

| Epoch | Block Number (Decimal) | Block Number (Hex) |
| ----- | ---------------------- | ------------------ |
| 0     | 0                      | 0x0                |
| 1     | 50                     | 0x32               |
| 2     | 100                    | 0x64               |
| 3     | 150                    | 0x96               |

## 💻 Commands

Now, we need enter `state-expiry-devnet` to run testing.

```bash
cd state-expiry-devnet
```

### Run state expiry client

All setup scripts are available in `state-expiry-devnet/scripts`. By default, it will create 3 nodes config, and start them.

#### Deploy 2 full nodes with state expiry, 1 normal full node 

```bash
bash scripts/test_deploy_nodes.sh start
```

#### Stop cluster nodes

```bash
bash scripts/test_deploy_nodes.sh stop
```

### Deploy ERC20 & Testing

Then enter `state-expiry-devnet/test-contract/deploy-token` to run all ERC20 scripts.

It mocks a fake ERC20 token to test contract slots expiry scenarios.

```bash
cd test-contract/deploy-token
```

#### Deploy ERC20 Token Contract

```bash
# install dependencies
npm install
# deploy BEP20 Token
npx hardhat run scripts/deploy.js
```

#### Transfer ERC20 Token & Read Balance

```bash
# This will transfer ERC20 Token from scripts/asset/test_account.json first account
# This script only transfer once
npx hardhat run scripts/transfer.js
# read sender & receiver's balance
npx hardhat run scripts/read-balance.js
```

Now, you could use the default users to `transfer` tokens, but when you stop using on-chain interactions after a period of time, only `readBalance` will work because it doesn't trigger on-chain state access.

Using the default config, after about 20 minutes, you cannot directly use the script to `transfer` tokens, as you may get `Access expired state` error.

In this case, you need to revive state. For details, see `test test_bep20_witness_revive`.

### Other testcases

Then enter `state-expiry-devnet/test-script` to run all golang scripts.

```bash
cd test-script
```

#### ETH transfer Loop

```bash
# The script keeps sending transactions until ctrl-c
go run test_eth_transfer.go
```

#### ERC20 Transfer Loop
```bash
# The script keeps sending transactions until ctrl-c
go run test_erc20_transfer.go
```

#### ERC20 Transfer Loop with expiry
The following script will send transactions such that some of them accesses expired slots.
```bash
# The script keeps sending transactions until ctrl-c
go run test_erc20_revive.go
```

### Useful RPC Commands

```bash
# query block height
curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":83}' 127.0.0.1:8502

# query block by number
curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x3", true],"id":83}' 127.0.0.1:8502

# query tx by hash
curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionByHash","params":["0x12beecfb1adb7d874c4714a7871e23cf70baef612235d1276568611460927f18"],"id":83}' 127.0.0.1:8502

# query tx receipt
curl -H "Content-Type: application/json" -X POST --data '{"jsonrpc":"2.0","method":"eth_getTransactionReceipt","params":["0x782192568c8ee3393e3f3e9b7ac46e231d3cbe0b96941b642e28220ba343209b"],"id":83}' 127.0.0.1:8502
```