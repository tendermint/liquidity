[![codecov](https://codecov.io/gh/tendermint/liquidity/branch/develop/graph/badge.svg)](https://codecov.io/gh/tendermint/liquidity?branch=develop)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tendermint/liquidity)](https://pkg.go.dev/github.com/tendermint/liquidity)

# Liquidity Module
the Liquidity module of the Cosmos-SDK, which serves AMM(Automated Market Makers) style decentralized liquidity providing and coin swap functions.

The module enable anyone to create a liquidity pool, deposit or withdraw coins from the liquidity pool, and request coin swap to the liquidity pool.

This module will be used in the Cosmos Hub, and any other blockchain based on Cosmos-SDK.

- The Cosmos Hub AMM should have strong philosophy of inclusiveness of users from different blockchains because its prime utility is inter-blockchain communication
- To possess such characteristics, the Liquidity module should provide most convenient ways for external users to come in and use the services provided by the Cosmos Hub
- The Liquidity module should not anticipate specific assets, such as Atom, into the process of user-flow in a forced manner. It is repeatedly proved that unnatural anticipation of native coin at unavoidable parts of process resulting in poor user attraction

## Key features

![new-amm-model](doc/img/new-amm-model.png)

**Combination of Traditional Orderbook-based Model and New AMM Model**

- Although new AMM model has multiple advantages over orderbook-based model,
combination of both models will create more enriched utilities for wider potential users
- We re-define the concept of a “swap order” in AMM as a “limit order with short lifetime”
in an orderbook-based exchange. Then, two concepts from two different models can be
combined as one united model so that the function can provide both ways to participate
into the trading and liquidity providing activities
- Although our first version of the Liquidity module will not provide limit order option, but
the base structure of the codebase is already anticipating such feature expansion in
near future
- Advantages of combined model
    - More freedom on how to provide liquidity : Limit orders
    - Combination of pool liquidity and limit order liquidity provides users more enriched trading environment

For detailed Mechanism, you can find on our recent [Paper](https://github.com/tendermint/liquidity/raw/develop/doc/Liquidity%20Module%20V1%20-%20Mechanism%20Explained.pdf)

## Installation

### Requirements

| Requirement | Notes            |
| ----------- | ---------------- |
| Go version  | Go1.15 or higher |
| Cosmos-SDK  | v0.40.0          |

### Get Liquidity Module source code 
```bash 
$ git clone https://github.com/tendermint/liquidity.git
$ cd liquidity
$ go mod tidy
```

### Build

```bash 
$ make build 
```
You can find the `liquidityd` binary on `build/`

### Install
```bash 
$ make install 
```

## liquidityd

### Tx


`$ ./liquidityd tx liquidity --help`     

```bash
Liquidity transaction subcommands

Usage:
  liquidityd tx liquidity [flags]
  liquidityd tx liquidity [command]

Available Commands:
  create-pool Create Liquidity pool with the specified pool-type, deposit coins
  deposit     Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins
  swap        Swap offer to the Liquidity pool with the specified the pool info with offer-coin, order-price
  withdraw    Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool
```



### Query

`$ liquidityd query liquidity --help`

```bash
Querying commands for the liquidity module

Usage:
  liquidityd query liquidity [flags]
  liquidityd query liquidity [command]

Available Commands:
  batch       Query details of a liquidity pool batch of the pool
  batches     Query for all liquidity pools batch
  deposits    Query for all deposit messages on the batch of the liquidity pool
  params      Query the current liquidity parameters information
  pool        Query details of a liquidity pool
  pools       Query for all liquidity pools
  swaps       Query for all swap messages on the batch of the liquidity pool
  withdraws   Query for all withdraw messages on the batch of the liquidity pool
```

#### A detailed document on client can be found here. [client.md](doc/client.md)

## Development

### Test
```bash 
$ make test
```

### Protobuf, Swagger

you can check local swagger doc page on `YOUR_API_SERVER(ex:127.0.0.1:1317)/swagger-liquidity/` if set `swagger = true` from `app.toml`
or see on [public swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/2.0.1)

generate `*.pb.go`, `*.pb.gw.go` files from `proto/*.proto`

```bash
$ make proto-gen
```
 
generate `swagger.yaml` from `proto/*.proto`

```bash
$ make proto-swagger-gen
```

## Setup local Testnet

Example of setup local testnet with test validator, user account

```bash
make install
liquidityd init testing --chain-id testing
liquidityd keys add validator --keyring-backend test
liquidityd keys add user1 --keyring-backend test
liquidityd add-genesis-account $(liquidityd keys show validator --keyring-backend test -a) 1000000000stake,1000000000token
liquidityd add-genesis-account $(liquidityd keys show user1 --keyring-backend test -a) 1000000000stake,1000000000atom
liquidityd gentx validator 1000000000stake --chain-id testing --keyring-backend test
liquidityd collect-gentxs
liquidityd start
```

### Broadcasting Txs with cli

Example of creating test liquidity pool 1 using cli

```bash
liquidityd tx liquidity create-pool 1 100000000stake,100000000token --from validator --keyring-backend test --chain-id testing -y
```

Example of creating test liquidity pool 2 using cli

```bash
liquidityd tx liquidity create-pool 1 100000000stake,100000000atom --from user1 --keyring-backend test --chain-id testing -y
```

Example of Swap request using cli

```bash
liquidityd tx liquidity swap 2 1 1000stake atom 1.15 --from validator --chain-id testing --keyring-backend test -y
```

### Broadcasting Txs with REST

Example of broadcast txs the new REST endpoint (via gRPC-gateway),

example of generating unsigned tx 

```bash
validator=$(liquidityd keys show validator --keyring-backend test -a)
liquidityd tx liquidity swap 2 1 1000stake atom 1.15 --from $validator --chain-id testing --generate-only > tx_swap.json
cat tx_swap.json
```
 
example of signing unsigned tx

```bash
liquidityd tx sign tx_swap.json --from validator --chain-id testing --keyring-backend test -y > tx_swap_signed.json
cat tx_swap_signed.json
```

example of encoding signed tx

```bash
liquidityd tx encode tx_swap_signed.json
```

example of the output: `CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAcSBBDAmgwaQHv+S5sy4L6I1/lLebkqiKeDX0xr6aBwd/XxaGvwTqUVTKsD7TqyBIeBxepg2a/7xi45Ws9xr8LLt7nAaR9DsdU=`


example of broadcasting txs using the [new REST endpoint (via gRPC-gateway, beta1)](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints)
need to enable API server for test

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```

## Export, Genesis State

### export empty state case
`liquidityd testnet --v 1` 

`liquidityd start --home ./mytestnet/node0/liquidityd/`

`liquidityd export  --home ./mytestnet/node0/liquidityd/`

```json
...
"liquidity": {
      "liquidity_pool_records": [],
      "params": {
        "init_pool_coin_mint_amount": "1000000",
        "liquidity_pool_creation_fee": [
          {
            "amount": "100000000",
            "denom": "stake"
          }
        ],
        "liquidity_pool_types": [
          {
            "description": "",
            "max_reserve_coin_num": 2,
            "min_reserve_coin_num": 2,
            "name": "DefaultPoolType",
            "pool_type_index": 1
          }
        ],
        "min_init_deposit_to_pool": "1000000",
        "swap_fee_rate": "0.003000000000000000"
      }
    },
    "mint": {
      "minter": {
        "annual_provisions": "130000037.646079971921585420",
        "inflation": "0.130000035046079271"
      },
      "params": {
        "blocks_per_year": "6311520",
        "goal_bonded": "0.670000000000000000",
        "inflation_max": "0.200000000000000000",
        "inflation_min": "0.070000000000000000",
        "inflation_rate_change": "0.130000000000000000",
        "mint_denom": "stake"
      }
    },

...
```

### pool created state export case

`liquidityd testnet --v 1`

`liquidityd start --home ./mytestnet/node0/liquidityd/`

`cat mytestnet/node0/liquidityd/config/genesis.json | grep chain_id`

`liquidityd tx liquidity create-pool 1 100000000reservecoin1,100000000reservecoin2 --from node0  --home ./mytestnet/node0/liquidityd/ --fees 2stake --chain-id <CHAIN-ID>`

`liquidityd export --home ./mytestnet/node0/liquidityd/`

```json
...
"liquidity": {
      "liquidity_pool_records": [
        {
          "batch_pool_deposit_msgs": [],
          "batch_pool_swap_msg_records": [],
          "batch_pool_swap_msgs": [],
          "batch_pool_withdraw_msgs": [],
          "liquidity_pool": {
            "pool_coin_denom": "cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs",
            "pool_id": "1",
            "pool_type_index": 1,
            "reserve_account_address": "cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs",
            "reserve_coin_denoms": [
              "reservecoin1",
              "reservecoin2"
            ]
          },
          "liquidity_pool_batch": {
            "batch_index": "4",
            "begin_height": "12",
            "deposit_msg_index": "1",
            "executed": true,
            "pool_id": "1",
            "swap_msg_index": "1",
            "withdraw_msg_index": "1"
          },
          "liquidity_pool_meta_data": {
            "pool_coin_total_supply": {
              "amount": "1000000",
              "denom": "cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs"
            },
            "pool_id": "1",
            "reserve_coins": [
              {
                "amount": "100000000",
                "denom": "reservecoin1"
              },
              {
                "amount": "100000000",
                "denom": "reservecoin2"
              }
            ]
          }
        }
      ],
      "params": {
        "init_pool_coin_mint_amount": "1000000",
        "liquidity_pool_creation_fee": [
          {
            "amount": "100000000",
            "denom": "stake"
          }
        ],
        "liquidity_pool_types": [
          {
            "description": "",
            "max_reserve_coin_num": 2,
            "min_reserve_coin_num": 2,
            "name": "DefaultPoolType",
            "pool_type_index": 1
          }
        ],
        "min_init_deposit_to_pool": "1000000",
        "swap_fee_rate": "0.003000000000000000"
      }
    },
...
```
 
## Resources
 - [Spec](x/liquidity/spec)
 - [Liquidity Module V1 Mechanism Paper](doc/Liquidity%20Module%20V1%20-%20Mechanism%20Explained.pdf)
 - [Proposal and milestone](https://github.com/b-harvest/Liquidity-Module-For-the-Hub)
 - [swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/2.0.1)
 - [godoc](https://pkg.go.dev/github.com/tendermint/liquidity)
 - [liquidityd client doc](doc/client.md)
 
