
# Liquidityd

Implemented tx cli

- [x]  `create-pool`   Create liquidity pool and deposit coins
- [x]  `deposit`       Deposit coins to a liquidity pool
- [x]  `swap`          Swap offer coin with demand coin from the liquidity pool with the given order price
- [x]  `withdraw`      Withdraw pool coin from the specified liquidity pool 

Implemented query cli 

- [x]    `batch`       Query details of a liquidity pool batch
- [x]    `deposit`     Query the deposit messages on the liquidity pool batch
- [x]    `deposits`    Query all deposit messages of the liquidity pool batch
- [x]    `params`      Query the values set as liquidity parameters
- [x]    `pool`        Query details of a liquidity pool
- [x]    `pools`       Query for all liquidity pools
- [x]    `swap`        Query for the swap message on the batch of the liquidity pool specified pool-id and msg-index
- [x]    `swaps`       Query all swap messages in the liquidity pool batch
- [x]    `withdraw`    Query the withdraw messages in the liquidity pool batch
- [x]    `withdraws`   Query for all withdraw messages on the liquidity pool batch

Progress REST/API

- [x] liquidity query endpoints of REST api using grpc model
- [x] broadcast txs using the new REST endpoint (via gRPC-gateway, beta1)


## Tx

`$ liquidityd tx liquidity --help`

```bash
Liquidity transaction subcommands

Usage:
  liquidityd tx liquidity [flags]
  liquidityd tx liquidity [command]

Available Commands:
  create-pool Create liquidity pool and deposit coins
  deposit     Deposit coins to a liquidity pool
  swap        Swap offer coin with demand coin from the liquidity pool with the given order price
  withdraw    Withdraw pool coin from the specified liquidity pool
```

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### tx create-pool

`$ liquidityd tx liquidity create-pool --help`

```bash
Create liquidity pool and deposit coins.

Example:
$ liquidityd tx liquidity create-pool 1 1000000000uatom,50000000000uusd --from mykey

This example creates a liquidity pool of pool-type 1 (two coins) and deposits 1000000000uatom and 50000000000uusd.
New liquidity pools can be created only for coin combinations that do not already exist in the network.

[pool-type]: The id of the liquidity pool-type. The only supported pool type is 1
[deposit-coins]: The amount of coins to deposit to the liquidity pool. The number of deposit coins must be 2 in pool type 1.

Usage:
  liquidityd tx liquidity create-pool [pool-type] [deposit-coins] [flags]
```

example tx command with result 

`$ liquidityd tx liquidity create-pool 1 1000000000uatom,50000000000uusd --from user1 --keyring-backend test --chain-id testing -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgCreatePool",
        "pool_creator_address": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft",
        "pool_type_id": 1,
        "deposit_coins": [
          {
            "denom": "uatom",
            "amount": "1000000000"
          },
          {
            "denom": "uusd",
            "amount": "50000000000"
          }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

result

```json
{
  "height": "5",
  "txhash": "C326C06CFB50589F72CBACD6F0028EE00B94F259C869D55653CEE11208531496",
  "codespace": "",
  "code": 0,
  "data": "0A0D0A0B6372656174655F706F6F6C",
  "raw_log": "...",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "create_pool",
          "attributes": [
            {
              "key": "pool_id",
              "value": "1"
            },
            {
              "key": "pool_type_id",
              "value": "1"
            },
            {
              "key": "pool_name",
              "value": "uatom/uusd/1"
            },
            {
              "key": "reserve_account",
              "value": "cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we"
            },
            {
              "key": "deposit_coins",
              "value": "1000000000uatom,50000000000uusd"
            },
            {
              "key": "pool_coin_denom",
              "value": "pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "create_pool"
            },
            {
              "key": "sender",
              "value": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft"
            },
            {
              "key": "sender",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft"
            },
            {
              "key": "module",
              "value": "liquidity"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we"
            },
            {
              "key": "amount",
              "value": "1000000000uatom,50000000000uusd"
            },
            {
              "key": "recipient",
              "value": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft"
            },
            {
              "key": "amount",
              "value": "1000000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
            },
            {
              "key": "recipient",
              "value": "cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl"
            },
            {
              "key": "sender",
              "value": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft"
            },
            {
              "key": "amount",
              "value": "100000000stake"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "163716",
  "tx": null,
  "timestamp": ""
}
```

already exist case, when duplicated request for same create pool

```json
{
  "height": "21",
  "txhash": "3077049475293028D4573197C84B4F1B4168D05532BD5BF3E97174B1A2BF520C",
  "codespace": "liquidity",
  "code": 11,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: the pool already exists",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "49392",
  "tx": null,
  "timestamp": ""
}
```

### tx deposit

`$ liquidityd tx liquidity deposit --help`

```bash 
Deposit coins a liquidity pool.

This deposit request is not processed immediately since it is accumulated in the liquidity pool batch.
All requests in a batch are treated equally and executed at the same swap price.

Example:
$ liquidityd tx liquidity deposit 1 100000000uatom,5000000000uusd --from mykey

This example request deposits 100000000uatom and 5000000000uusd to pool-id 1.
Deposits must be the same coin denoms as the reserve coins.

[pool-id]: The pool id of the liquidity pool
[deposit-coins]: The amount of coins to deposit to the liquidity pool

Usage:
  liquidityd tx liquidity deposit [pool-id] [deposit-coins] [flags]
```

example tx command with result 


`$ liquidityd tx liquidity deposit 1 100000000uatom,5000000000uusd --from validator --keyring-backend test --chain-id testing -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgDepositWithinBatch",
        "depositor_address": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3",
        "pool_id": "1",
        "deposit_coins": [
          {
            "denom": "uatom",
            "amount": "100000000"
          },
          {
            "denom": "uusd",
            "amount": "5000000000"
          }
        ]
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

result 

```
{
  "height": "458",
  "txhash": "8D8FA31125AB2A984D28F362ADC05946208C0E7927B13F984D9AD6E8E5327782",
  "codespace": "",
  "code": 0,
  "data": "0A160A146465706F7369745F77697468696E5F6261746368",
  "raw_log": "...",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "deposit_within_batch",
          "attributes": [
            {
              "key": "pool_id",
              "value": "1"
            },
            {
              "key": "batch_index",
              "value": "1"
            },
            {
              "key": "msg_index",
              "value": "1"
            },
            {
              "key": "deposit_coins",
              "value": "100000000uatom,5000000000uusd"
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "deposit_within_batch"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "module",
              "value": "liquidity"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "100000000uatom,5000000000uusd"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "79385",
  "tx": null,
  "timestamp": ""
}
```


### tx swap

`$ liquidityd tx liquidity swap --help`

```bash  
Swap offer coin with demand coin from the liquidity pool with the given order price.

This swap request is not processed immediately since it is accumulated in the liquidity pool batch.
All requests in a batch are treated equally and executed at the same swap price.
The order of swap requests is ignored since the universal swap price is calculated in every batch to prevent front running.

The requested swap is executed with a swap price that is calculated from the given swap price function of the pool, the other swap requests, and the liquidity pool coin reserve status.
Swap orders are executed only when the execution swap price is equal to or greater than the submitted order price of the swap order.

Example:
$ liquidity liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from mykey

For this example, imagine that an existing liquidity pool has with 1000000000uatom and 50000000000uusd.
This example request swaps 50000000uusd for at least 950000uatom with the order price of 0.019 and swap fee rate of 0.003.
A sufficient balance of half of the swap-fee-rate of the offer coin is required to reserve the offer coin fee.

The order price is the exchange ratio of X/Y, where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically.
Increasing order price reduces the possibility for your request to be processed and results in buying uatom at a lower price than the pool price.

For explicit calculations, The swap fee rate must be the value that set as liquidity parameter in the current network.
The only supported swap-type is 1. For the detailed swap algorithm, see https://github.com/tendermint/liquidity

[pool-id]: The pool id of the liquidity pool
[swap-type]: The swap type of the swap message. The only supported swap type is 1 (instant swap).
[offer-coin]: The amount of offer coin to swap
[demand-coin-denom]: The denomination of the coin to exchange with offer coin
[order-price]: The limit order price for the swap order. The price is the exchange ratio of X/Y where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically
[swap-fee-rate]: The swap fee rate to pay for swap that is proportional to swap amount. The swap fee rate must be the value that set as liquidity parameter in the current network.

Usage:
  liquidityd tx liquidity swap [pool-id] [swap-type] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate] [flags]
```

example tx command with result 

`$ liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from validator --chain-id testing --keyring-backend test -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgSwapWithinBatch",
        "swap_requester_address": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3",
        "pool_id": "1",
        "swap_type_id": 1,
        "offer_coin": {
          "denom": "uusd",
          "amount": "50000000"
        },
        "demand_coin_denom": "uatom",
        "offer_coin_fee": {
          "denom": "uusd",
          "amount": "75000"
        },
        "order_price": "0.019000000000000000"
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

result 

```
{
  "height": "178",
  "txhash": "AA9A3A50D9AC639730F61824AA2BD3BA9EBCCEA7E52147353C0E680041F21243",
  "codespace": "",
  "code": 0,
  "data": "0A130A11737761705F77697468696E5F6261746368",
  "raw_log": "...",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "swap_within_batch"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "module",
              "value": "liquidity"
            }
          ]
        },
        {
          "type": "swap_within_batch",
          "attributes": [
            {
              "key": "pool_id",
              "value": "1"
            },
            {
              "key": "batch_index",
              "value": "1"
            },
            {
              "key": "msg_index",
              "value": "1"
            },
            {
              "key": "swap_type_id",
              "value": "1"
            },
            {
              "key": "offer_coin_denom",
              "value": "uusd"
            },
            {
              "key": "offer_coin_amount",
              "value": "50000000"
            },
            {
              "key": "offer_coin_fee_amount",
              "value": "75000"
            },
            {
              "key": "demand_coin_denom",
              "value": "uatom"
            },
            {
              "key": "order_price",
              "value": "0.019000000000000000"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "50000000uusd"
            },
            {
              "key": "recipient",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "75000uusd"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "95327",
  "tx": null,
  "timestamp": ""
}
```

### tx withdraw

`$ liquidityd tx liquidity withdraw --help`

```bash 
Withdraw pool coin from the specified liquidity pool.

This swap request is not processed immediately since it is accumulated in the liquidity pool batch.
All requests in a batch are treated equally and executed at the same swap price.

Example:
$ liquidityd tx liquidity withdraw 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from mykey

This example request withdraws 10000 pool coin from the specified liquidity pool.
The appropriate pool coin must be requested from the specified pool.

[pool-id]: The pool id of the liquidity pool
[pool-coin]: The amount of pool coin to withdraw from the liquidity pool

Usage:
  liquidityd tx liquidity withdraw [pool-id] [pool-coin] [flags]
```

check the balance before withdraw

`$ liquidityd query bank balances cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3`

```
balances:
- amount: "99899"
  denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
- amount: "1000000"
  denom: poolA4648A10F8D43B8EE4D915A35CB292618215D9F60CE3E2E29216489CF1FAE049
- amount: "8890000000"
  denom: stake
- amount: "9901196107"
  denom: uatom
- amount: "494939925000"
  denom: uusd
pagination:
  next_key: null
  total: "0"
```

example tx command with result

`$ liquidityd tx liquidity withdraw 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from validator --chain-id testing --keyring-backend test -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgWithdrawWithinBatch",
        "withdrawer_address": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3",
        "pool_id": "1",
        "pool_coin": {
          "denom": "pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295",
          "amount": "10000"
        }
      }
    ],
    "memo": "",
    "timeout_height": "0",
    "extension_options": [],
    "non_critical_extension_options": []
  },
  "auth_info": {
    "signer_infos": [],
    "fee": {
      "amount": [],
      "gas_limit": "200000",
      "payer": "",
      "granter": ""
    }
  },
  "signatures": []
}
```

result 

```
{
  "height": "562",
  "txhash": "BE8827F69E8BC5909A0FFC713B6D267606A91A1CFA07552E69020638E9E1D563",
  "codespace": "",
  "code": 0,
  "data": "0A170A1577697468647261775F77697468696E5F6261746368",
  "raw_log": "...",
  "logs": [
    {
      "msg_index": 0,
      "log": "",
      "events": [
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "withdraw_within_batch"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "module",
              "value": "liquidity"
            }
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
            }
          ]
        },
        {
          "type": "withdraw_within_batch",
          "attributes": [
            {
              "key": "pool_id",
              "value": "1"
            },
            {
              "key": "batch_index",
              "value": "2"
            },
            {
              "key": "msg_index",
              "value": "1"
            },
            {
              "key": "pool_coin_denom",
              "value": "pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
            },
            {
              "key": "pool_coin_amount",
              "value": "10000"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "67701",
  "tx": null,
  "timestamp": ""
}
```

balances after withdraw

`$ liquidityd query bank balances cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3`

```
balances:
- amount: "89899"
  denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
- amount: "1000000"
  denom: poolA4648A10F8D43B8EE4D915A35CB292618215D9F60CE3E2E29216489CF1FAE049
- amount: "8890000000"
  denom: stake
- amount: "9911156180"
  denom: uatom
- amount: "495438924678"
  denom: uusd
pagination:
  next_key: null
  total: "0"
```

## Query

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

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### query batch
`$ liquidityd query liquidity batch --help`
```bash
Query details of a liquidity pool batch
Example:
$ liquidityd query liquidity batch 1

Usage:
  liquidityd query liquidity batch [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity batch 1`

```bash
batch:
  begin_height: "563"
  deposit_msg_index: "2"
  executed: false
  index: "3"
  pool_id: "1"
  swap_msg_index: "2"
  withdraw_msg_index: "2"
```

### query deposits
`$ liquidityd query liquidity deposits --help`
```bash
Query all deposit messages of the liquidity pool batch on the specified pool

If batch messages are normally processed from the endblock, the resulting state is applied and the messages are removed in the beginning of next block.
To query for past blocks, query the block height using the REST/gRPC API of a node that is not pruned.

Example:
$ liquidityd query liquidity deposits 1

Usage:
  liquidityd query liquidity deposits [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity deposits 1`

```bash
deposits:
- executed: true
  msg:
    deposit_coins:
    - amount: "100000000"
      denom: uatom
    - amount: "5000000000"
      denom: uusd
    depositor_address: cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3
    pool_id: "1"
  msg_height: "458"
  msg_index: "1"
  succeeded: true
  to_be_deleted: true
pagination:
  next_key: null
  total: "1"
```

empty case

`$ liquidityd query liquidity deposits 1`

```bash 
deposits: []
pagination:
  next_key: null
  total: "0"
```

### query pool

`$ liquidityd query liquidity pool  --help`

```   
Query details of a liquidity pool
Example:
$ liquidityd query liquidity pool 1

Usage:
  liquidityd query liquidity pool [pool-id] [flags]
```

example query command with result 

`$ liquidityd query liquidity pool 1`
 
```bash
pool:
  id: "1"
  pool_coin_denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
  reserve_account_address: cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
  reserve_coin_denoms:
  - uatom
  - uusd
  type_id: 1
```

example query reserve coins of the pool

`$ liquidityd query bank balances cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we`

```bash
balances:
- amount: "999003494"
  denom: uatom
- amount: "50050075000"
  denom: uusd
pagination:
  next_key: null
  total: "0"
```


example query total supply the pool coin

`$ liquidityd query bank total --denom=pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295`

```bash
amount: "1000000"
denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
```

### query pools
`$ liquidityd query liquidity pools  --help`
```   
Query details about all liquidity pools on a network.
Example:
$ liquidityd query liquidity pools

Usage:
  liquidityd query liquidity pools [flags]
```

example query command with result 

`$ liquidityd query liquidity pools`
 
```bash
pagination:
  next_key: null
  total: "2"
pools:
- id: "1"
  pool_coin_denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
  reserve_account_address: cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
  reserve_coin_denoms:
  - uatom
  - uusd
  type_id: 1
- id: "2"
  pool_coin_denom: poolA4648A10F8D43B8EE4D915A35CB292618215D9F60CE3E2E29216489CF1FAE049
  reserve_account_address: cosmos153jg5y8c6sacaexezk34ev5jvxpptk0kscrx0x
  reserve_coin_denoms:
  - stake
  - uusd
  type_id: 1
```

### query params

example query command with result 

`$ liquidityd query liquidity params`

```bash
init_pool_coin_mint_amount: "1000000"
max_order_amount_ratio: "0.100000000000000000"
max_reserve_coin_amount: "0"
min_init_deposit_amount: "1000000"
pool_creation_fee:
- amount: "100000000"
  denom: stake
pool_types:
- description: ""
  id: 1
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: DefaultPoolType
swap_fee_rate: "0.003000000000000000"
unit_batch_height: 1
withdraw_fee_rate: "0.003000000000000000"
```


### query swaps
`$ liquidityd query liquidity swaps --help`
```bash
Query all swap messages in the liquidity pool batch for the specified pool-id

If batch messages are normally processed from the endblock,
the resulting state is applied and the messages are removed in the beginning of next block.
To query for past blocks, query the block height using the REST/gRPC API of a node that is not pruned.

Example:
$ liquidityd query liquidity swaps 1

Usage:
  liquidityd query liquidity swaps [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity swaps 1`
```bash
pagination:
  next_key: null
  total: "1"
swaps:
- exchanged_offer_coin:
    amount: "50000000"
    denom: uusd
  executed: true
  msg:
    demand_coin_denom: uatom
    offer_coin:
      amount: "50000000"
      denom: uusd
    offer_coin_fee:
      amount: "75000"
      denom: uusd
    order_price: "0.019000000000000000"
    pool_id: "1"
    swap_requester_address: cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3
    swap_type_id: 1
  msg_height: "178"
  msg_index: "1"
  order_expiry_height: "178"
  remaining_offer_coin:
    amount: "0"
    denom: uusd
  reserved_offer_coin_fee:
    amount: "0"
    denom: uusd
  succeeded: true
  to_be_deleted: true
```

empty case

`$ liquidityd query liquidity swaps 1`
```bash 
pagination:
  next_key: null
  total: "0"
swaps: []
```


### query withdraws
`$ liquidityd query liquidity withdraws --help`
```bash
Query all withdraw messages on the liquidity pool batch for the specified pool-id

If batch messages are normally processed from the endblock,
the resulting state is applied and the messages are removed in the beginning of next block.
To query for past blocks, query the block height using the REST/gRPC API of a node that is not pruned.

Example:
$ liquidityd query liquidity withdraws 1

Usage:
  liquidityd query liquidity withdraws [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity withdraws 1`

```bash
pagination:
  next_key: null
  total: "1"
withdraws:
- executed: true
  msg:
    pool_coin:
      amount: "10000"
      denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
    pool_id: "1"
    withdrawer_address: cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3
  msg_height: "562"
  msg_index: "1"
  succeeded: true
  to_be_deleted: true
```

empty case
`$ liquidityd query liquidity withdraws 1`
```bash 
pagination:
  next_key: null
  total: "0"
withdraws: []
```

## REST/API

You can check local swagger doc page on `YOUR_API_SERVER(ex:127.0.0.1:1317)/swagger-liquidity/` if set `swagger = true` from `app.toml`
or see on [public swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs)

According to [migrating-to-new-rest-endpoints](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints), the POST endpoints of the New gGPC-gateway REST are N/A and guided directly to use Protobuf, need to use `cli` or `localhost:1317/cosmos/tx/v1beta1/txs` for broadcast txs temporarily

example of broadcasting txs using the [new REST endpoint (via gRPC-gateway, beta1)](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints)

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```
