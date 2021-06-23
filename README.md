[![codecov](https://codecov.io/gh/tendermint/liquidity/branch/develop/graph/badge.svg)](https://codecov.io/gh/tendermint/liquidity?branch=develop)
[![PkgGoDev](https://pkg.go.dev/badge/github.com/tendermint/liquidity)](https://pkg.go.dev/github.com/tendermint/liquidity)

# Liquidity Module

The liquidity module serves Automated Market Maker (AMM)-style decentralized liquidity by providing liquidity activities and coin swap functions.

The module enables users to create a liquidity pool, make deposits and withdrawals, and request coin swaps from the liquidity pool.

This module can be used in the [Cosmos Hub](https://hub.cosmos.network/main/hub-overview/overview.html) and any other [Cosmos SDK](https://cosmos.network/)-based blockchain projects.

- The Cosmos Hub AMM applies a strong philosophy of inclusiveness for users from different blockchains with its prime utility of inter-blockchain communication.
- To achieve heterogeneous blockchain adoption, the liquidity module provides convenient entry points for external users to come in and use the services that are provided by the Cosmos Hub.
- The liquidity module does not anticipate specific assets, such as ATOM, into the user workflow. Data shows that unnatural anticipation of native coin at unavoidable parts of the process results in poor user attraction.

## Key features

![new-amm-model](doc/img/new-amm-model.png)

**Combination of traditional orderbook-based model and new AMM model**

- With multiple advantages over order book-based models, the liquidity module combines a batch-based order book matching algorithm with AMM to create enriched utilities for more potential users.
- The liquidity module redefines the concept of a “swap order” in AMM as a “limit order with a short lifetime” in an order book-based exchange. By combining these concepts from two different models as one united model, the function supports both ways to participate in trading and liquidity-providing activities.
- Limit order options are not supported in the first version of the liquidity module, but the base structure of the codebase anticipates and supports feature expansion.
- Advantages of the combined model
    - More freedom on ways to provide liquidity, planned expansion for limit orders
    - The combination of pool liquidity and limit order liquidity provide users with a more enriched trading environment

For details, see the [Liquidity Module Light Paper](doc/LiquidityModuleLightPaper_EN.pdf).

## Installation

### Requirements

Requirement | Notes
----------- | -----------------
Go version  | Go1.15 or higher
Cosmos SDK  | v0.42.4 or higher

### Get Liquidity Module source code

```bash
$ git clone https://github.com/tendermint/liquidity.git
$ cd liquidity
$ go mod tidy
```

### Build

```bash
# The `liquidityd` binary is in the build directory.
$ make build
```

### Install

```bash
$ make install
```

## Usage of CLI Commands

With the exception of creating the liquidity pool, all commands are implemented to execute on the batch.

### Transactions

`$ liquidityd tx liquidity --help`

```bash
Liquidity transaction subcommands

Usage:
  liquidityd tx liquidity [flags]
  liquidityd tx liquidity [command]

Available Commands:
  create-pool Create liquidity pool and deposit coins
  deposit     Deposit coins to a liquidity pool
  swap        Swap offer coin with demand coin
  withdraw    Withdraw pool coin
```

### Queries

`$ liquidityd query liquidity --help`

```bash
Querying commands for the liquidity module

Usage:
  liquidityd query liquidity [flags]
  liquidityd query liquidity [command]

Available Commands:
  batch       Query details of a liquidity pool batch
  deposit     Query the deposit messages on the liquidity pool batch
  deposits    Query all deposit messages of the liquidity pool batch
  params      Query the values set as liquidity parameters
  pool        Query details of a liquidity pool
  pools       Query for all liquidity pools
  swap        Query for the swap message on the batch of the liquidity pool specified pool-id and msg-index
  swaps       Query all swap messages in the liquidity pool batch
  withdraw    Query the withdraw messages in the liquidity pool batch
  withdraws   Query for all withdraw messages on the liquidity pool batch
```

#### A detailed document on client can be found here. [client.md](doc/client.md)

## Development

### Test

```bash
$ make test-all
```

### 1\. Setup local testnet using script

```bash
# This script bootstraps a single local testnet.
# Note that config, data, and keys are created in the ./data/localnet folder and
# RPC, GRPC, and REST ports are all open.
$ make localnet
```

### 1.1 Broadcast transactions using CLI command-line interface

Sample scripts are provided in [scripts](https://github.com/tendermint/liquidity/tree/develop/scripts) folder to help you to test the liquidity module interface.

### 2. Manually set up a local testnet

```bash
# Build
make install

# Initialize and add keys
liquidityd init testing --chain-id testing
liquidityd keys add validator --keyring-backend test
liquidityd keys add user1 --keyring-backend test

# Add genesis accounts and provide coins to the accounts
liquidityd add-genesis-account $(liquidityd keys show validator --keyring-backend test -a) 10000000000stake,10000000000uatom,500000000000uusd
liquidityd add-genesis-account $(liquidityd keys show user1 --keyring-backend test -a) 10000000000stake,10000000000uatom,500000000000uusd

# Create gentx and collect
liquidityd gentx validator 1000000000stake --chain-id testing --keyring-backend test
liquidityd collect-gentxs

# Start
liquidityd start
```

### 2.1 Broadcast transactions using CLI commands

```bash
# An example of creating liquidity pool 1
liquidityd tx liquidity create-pool 1 1000000000uatom,50000000000uusd --from user1 --keyring-backend test --chain-id testing -y

# An example of creating liquidity pool 2
liquidityd tx liquidity create-pool 1 10000000stake,10000000uusd --from validator --keyring-backend test --chain-id testing -y

# An example of requesting swap
liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from validator --chain-id testing --keyring-backend test -y

# An example of generating unsigned tx
validator=$(liquidityd keys show validator --keyring-backend test -a)
liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from $validator --chain-id testing --generate-only > tx_swap.json
cat tx_swap.json

# Sign the unsigned tx
liquidityd tx sign tx_swap.json --from validator --chain-id testing --keyring-backend test -y > tx_swap_signed.json
cat tx_swap_signed.json

# Encode the signed tx
liquidityd tx encode tx_swap_signed.json
```

### 2.2 Broadcast transactions using REST APIs

For an example of broadcasting transactions using REST API (via gRPC-gateway), see Cosmos SDK [Migrating to New REST Endpoints](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints). Testing requires that the API server is enabled in `$HOME/.liquidityapp/config/app.toml`.

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"Cp0BCpoBCigvdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcFdpdGhpbkJhdGNoEm4KLWNvc21vczE4cWM2ZGwwNDZ1a3V0MjN3NnF1dndmenBmeWhncDJmeHFkcXAwNhACGAEiEAoEdXVzZBIINTAwMDAwMDAqBXVhdG9tMg0KBHV1c2QSBTc1MDAwOhExOTAwMDAwMDAwMDAwMDAwMBJYClAKRgofL2Nvc21vcy5jcnlwdG8uc2VjcDI1NmsxLlB1YktleRIjCiEDsouFptHWGniIBzFrsE26PcfH950qjnf4RaEsd+g2fA0SBAoCCH8YAxIEEMCaDBpAOI3k8fay9TziZbl+eNCqmPEF7tWXua3ad0ldNR6XOgZjKRBP9sQSxCtaRFnqc6Avep9C4Rjt+CHDahRNpZ8u3A==","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```

### 2.3 Export Genesis State

`$ liquidityd export`

### Export empty state case

```json
{
  "liquidity": {
    "params": {
      "circuit_breaker_enabled": false,
      "init_pool_coin_mint_amount": "1000000",
      "max_order_amount_ratio": "0.100000000000000000",
      "max_reserve_coin_amount": "0",
      "min_init_deposit_amount": "1000000",
      "pool_creation_fee": [
        {
          "amount": "40000000",
          "denom": "stake"
        }
      ],
      "pool_types": [
        {
          "description": "Standard liquidity pool with pool price function X/Y, ESPM constraint, and two kinds of reserve coins",
          "id": 1,
          "max_reserve_coin_num": 2,
          "min_reserve_coin_num": 2,
          "name": "StandardLiquidityPool"
        }
      ],
      "swap_fee_rate": "0.003000000000000000",
      "unit_batch_height": 1,
      "withdraw_fee_rate": "0.000000000000000000"
    },
    "pool_records": []
  }
}
```

### Export when some states exist

```json
{
  "liquidity": {
    "params": {
      "circuit_breaker_enabled": false,
      "init_pool_coin_mint_amount": "1000000",
      "max_order_amount_ratio": "0.100000000000000000",
      "max_reserve_coin_amount": "0",
      "min_init_deposit_amount": "1000000",
      "pool_creation_fee": [
        {
          "amount": "40000000",
          "denom": "stake"
        }
      ],
      "pool_types": [
        {
          "description": "Standard liquidity pool with pool price function X/Y, ESPM constraint, and two kinds of reserve coins",
          "id": 1,
          "max_reserve_coin_num": 2,
          "min_reserve_coin_num": 2,
          "name": "StandardLiquidityPool"
        }
      ],
      "swap_fee_rate": "0.003000000000000000",
      "unit_batch_height": 1,
      "withdraw_fee_rate": "0.000000000000000000"
    },
    "pool_records": [
      {
        "deposit_msg_states": [],
        "pool": {
          "id": "1",
          "pool_coin_denom": "pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295",
          "reserve_account_address": "cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we",
          "reserve_coin_denoms": [
            "uatom",
            "uusd"
          ],
          "type_id": 1
        },
        "pool_batch": {
          "begin_height": "563",
          "deposit_msg_index": "2",
          "executed": false,
          "index": "3",
          "pool_id": "1",
          "swap_msg_index": "2",
          "withdraw_msg_index": "2"
        },
        "pool_metadata": {
          "pool_coin_total_supply": {
            "amount": "1089899",
            "denom": "pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
          },
          "pool_id": "1",
          "reserve_coins": [
            {
              "amount": "1088843820",
              "denom": "uatom"
            },
            {
              "amount": "54551075322",
              "denom": "uusd"
            }
          ]
        },
        "swap_msg_states": [],
        "withdraw_msg_states": []
      },
      {
        "deposit_msg_states": [],
        "pool": {
          "id": "2",
          "pool_coin_denom": "poolA4648A10F8D43B8EE4D915A35CB292618215D9F60CE3E2E29216489CF1FAE049",
          "reserve_account_address": "cosmos153jg5y8c6sacaexezk34ev5jvxpptk0kscrx0x",
          "reserve_coin_denoms": [
            "stake",
            "uusd"
          ],
          "type_id": 1
        },
        "pool_batch": {
          "begin_height": "0",
          "deposit_msg_index": "1",
          "executed": false,
          "index": "1",
          "pool_id": "2",
          "swap_msg_index": "1",
          "withdraw_msg_index": "1"
        },
        "pool_metadata": {
          "pool_coin_total_supply": {
            "amount": "1000000",
            "denom": "poolA4648A10F8D43B8EE4D915A35CB292618215D9F60CE3E2E29216489CF1FAE049"
          },
          "pool_id": "2",
          "reserve_coins": [
            {
              "amount": "10000000",
              "denom": "stake"
            },
            {
              "amount": "10000000",
              "denom": "uusd"
            }
          ]
        },
        "swap_msg_states": [],
        "withdraw_msg_states": []
      }
    ]
  }
}
```

### Protobuf and Swagger

The API documentation for the liquidity module is available on `http://localhost:1317/swagger-liquidity/` after you successfully boostrap a testnet in your local computer.

You must set `swagger` config to `true` in `$HOME/.liquidityapp/config/app.toml`. The public Swagger API docs are also available on [Cosmos SDK Liquidity Module - REST and gRPC Gateway docs](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs).

```bash
# Generate `*.pb.go`, `*.pb.gw.go` files from `proto/*.proto`
$ make proto-gen

# Generate `swagger.yaml` from `proto/*.proto`
$ make proto-swagger-gen
```

## Resources

To learn more about the liquidity module, check out the following resources:

 - [Liquidity Module Spec](x/liquidity/spec)
 - [Liquidity Module Lite Paper (English)](doc/LiquidityModuleLightPaper_EN.pdf)
 - [Liquidity Module Lite Paper (Korean)](doc/LiquidityModuleLightPaper_KO.pdf)
 - [Liquidity Module Lite Paper (Chinese)](doc/LiquidityModuleLightPaper_ZH.pdf)
 - [Proposal and milestone](https://github.com/b-harvest/Liquidity-Module-For-the-Hub)
 - [Swagger HTTP API doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs)
 - [godoc](https://pkg.go.dev/github.com/tendermint/liquidity)
 - [Client doc](doc/client.md)
 - [Performance Testing](doc/Performance%20Testing%20for%20Liquidity%20Module.pdf)
 