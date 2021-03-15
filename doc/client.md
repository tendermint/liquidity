
# Liquidityd

Implemented tx cli

- [x]  `create-pool`   Create Liquidity pool with the specified pool-type, deposit coins
- [x]  `deposit`       Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins
- [x]  `swap`          Swap offer to the Liquidity pool with the specified the pool info with offer-coin, order-price
- [x]  `withdraw`      Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool 

Implemented query cli 

- [x]    `batch`       Query details of a liquidity pool batch of the pool
- [x]    `batches`     Query for all liquidity pools batch
- [x]    `deposits`    Query for all deposit messages on the batch of the liquidity pool
- [x]    `params`      Query the current liquidity parameters information
- [x]    `pool`        Query details of a liquidity pool
- [x]    `pools`       Query for all liquidity pools
- [x]    `swaps`       Query for all swap messages on the batch of the liquidity pool
- [x]    `withdraws`   Query for all withdraw messages on the batch of the liquidity pool

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
  create-pool Create Liquidity pool with the specified pool-type, deposit coins
  deposit     Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins
  swap        Swap offer to the Liquidity pool with the specified the pool info with offer-coin, order-price
  withdraw    Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool
```

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### tx create-pool

`$ liquidityd tx liquidity create-pool --help`

```bash
Create Liquidity pool with the specified pool-type-index, deposit coins for reserve

Example:
$ liquidity tx liquidity create-pool 1 100000000stake,100000000token --from mykey

Currently, only the default pool-type-index 1 is available on this version
the number of deposit coins must be two in the pool-type-index 1

{"pool_type_id":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}

Usage:
  liquidityd tx liquidity create-pool [pool-type-index] [deposit-coins] [flags]


```

example tx command with result 

`$ liquidityd tx liquidity create-pool 1 100000000stake,100000000token --from validator --keyring-backend test --chain-id testing -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgCreatePool",
        "pool_creator_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
        "pool_type_id": 1,
        "deposit_coins": [
          {
            "denom": "stake",
            "amount": "100000000"
          },
          {
            "denom": "token",
            "amount": "100000000"
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
  "height": "6",
  "txhash": "8E75D2210BC9C569ECFE53139803501AAF9ED24F567B96718E2464CBF1384E7F",
  "codespace": "",
  "code": 0,
  "data": "0A170A156372656174655F6C69717569646974795F706F6F6C",
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
              "key": "reserve_coin_denoms",
              "value": "stake/token/1"
            },
            {
              "key": "reserve_account",
              "value": "cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u"
            },
            {
              "key": "deposit_coins",
              "value": "100000000stake,100000000token"
            },
            {
              "key": "pool_coin_denom",
              "value": "pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07"
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
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "sender",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "sender",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
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
              "value": "cosmos18l9ktac2vf2qyf8a8hjahh47995ymknzg8my6t"
            },
            {
              "key": "amount",
              "value": "100000000stake"
            },
            {
              "key": "recipient",
              "value": "cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u"
            },
            {
              "key": "amount",
              "value": "100000000stake,100000000token"
            },
            {
              "key": "recipient",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "amount",
              "value": "1000000pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "156795",
  "tx": null,
  "timestamp": ""
}
```

already exist case, when duplicated request for same create pool

```json
{
  "height": "35",
  "txhash": "1A2740EDE76425E12E5600AC452A58B1CEDDF3FECD9BCF501C192C27EA5342E6",
  "codespace": "liquidity",
  "code": 11,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: the pool already exists",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "48408",
  "tx": null,
  "timestamp": ""
}
```

### tx deposit

`$ liquidityd tx liquidity deposit --help  `

```bash 
./liquidityd tx liquidity deposit --help 
Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins for reserve

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity deposit 1 100000000stake,100000000token --from mykey

You should deposit the same coin as the reserve coin.

Usage:
  liquidityd tx liquidity deposit [pool-id] [deposit-coins] [flags]

```

example tx command with result 


`$ liquidityd tx liquidity deposit 1 50000000stake,50000000token --from validator --keyring-backend test --chain-id testing -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgDepositWithinBatch",
        "depositor_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
        "pool_id": "1",
        "deposit_coins": [
          {
            "denom": "stake",
            "amount": "50000000"
          },
          {
            "denom": "token",
            "amount": "50000000"
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
  "height": "51",
  "txhash": "8361F43BE0A37785C0ADE807FEE088592A783C0A4234915C9E90FCD87F12F88A",
  "codespace": "",
  "code": 0,
  "data": "0A1B0A196465706F7369745F746F5F6C69717569646974795F706F6F6C",
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
              "value": "50000000stake,50000000token"
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "deposit_to_pool"
            },
            {
              "key": "sender",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
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
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "amount",
              "value": "50000000stake,50000000token"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "67334",
  "tx": null,
  "timestamp": ""
}
```


### tx swap

`$ liquidityd tx liquidity swap --help`

```bash  
swap [pool-id] [swap-type-index] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate]

Swap offer to the Liquidity pool with the specified pool-id, swap-type-index demand-coin-denom
with the coin and the price you're offering

this requests are stacked in the batch of the liquidity pool, not immediately processed and
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity swap 2 1 100000000stake token 0.9 0.003 --from mykey

You should request the same each field as the pool.

Must have sufficient balance half the of the swapFee Rate of the offer coin to reserve offer coin fee.

For explicit calculations, you must enter the params.swap_fee_rate value of the current parameter state.

Currently, only the default pool-type, swap-type-index 1 is available on this version
The detailed swap algorithm can be found here.
https://github.com/tendermint/liquidity
```

example tx command with result 

`$ liquidityd tx liquidity swap 1 1 1000token stake 0.9 0.003 --from validator --chain-id testing --keyring-backend test -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgSwapWithinBatch",
        "swap_requester_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
        "pool_id": "1",
        "swap_type": 1,
        "offer_coin": {
          "denom": "token",
          "amount": "1000"
        },
        "demand_coin_denom": "stake",
        "offer_coin_fee": {
          "denom": "token",
          "amount": "1"
        },
        "order_price": "0.900000000000000000"
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
  "height": "80",
  "txhash": "6C8E945F550DB638635A4EE5D68BF478012AAC47F2B7AD91FB702A59737B834C",
  "codespace": "",
  "code": 0,
  "data": "0A060A0473776170",
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
              "value": "swap"
            },
            {
              "key": "sender",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "sender",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
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
              "key": "swap_type",
              "value": "1"
            },
            {
              "key": "offer_coin_denom",
              "value": "token"
            },
            {
              "key": "offer_coin_amount",
              "value": "1000"
            },
            {
              "key": "offer_coin_fee_amount",
              "value": "1"
            },
            {
              "key": "demand_coin_denom",
              "value": "stake"
            },
            {
              "key": "order_price",
              "value": "0.900000000000000000"
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
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "amount",
              "value": "1000token"
            },
            {
              "key": "recipient",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "amount",
              "value": "1token"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "92558",
  "tx": null,
  "timestamp": ""
}
```

### tx withdraw

`$ liquidityd tx liquidity withdraw --help  `

```bash 
Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity withdraw 1 1000pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07 --from mykey

You should request the matched pool-coin as the pool.

Usage:
  liquidityd tx liquidity withdraw [pool-id] [pool-coin] [flags]
```

check the balance before broadcast tx

`$ liquidityd query bank balances cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098`

```
balances:
- amount: "1500000"
  denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
- amount: "750000999"
  denom: stake
- amount: "849998999"
  denom: token
```

example tx command with result

`$ liquidityd tx liquidity withdraw 1 500000pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07 --from validator --chain-id testing --keyring-backend test -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgWithdrawWithinBatch",
        "withdrawer_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
        "pool_id": "1",
        "pool_coin": {
          "denom": "pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07",
          "amount": "500000"
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
  "height": "220",
  "txhash": "187AA85E64062C10DB0FD9102B37307E364C047FBEB1B0A8D43826E8A3E687EC",
  "codespace": "",
  "code": 0,
  "data": "0A1E0A1C77697468647261775F66726F6D5F6C69717569646974795F706F6F6C",
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
              "value": "withdraw_from_pool"
            },
            {
              "key": "sender",
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
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
              "value": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098"
            },
            {
              "key": "amount",
              "value": "500000pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07"
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
              "value": "pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07"
            },
            {
              "key": "pool_coin_amount",
              "value": "500000"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "67718",
  "tx": null,
  "timestamp": ""
}
```

balances after withdraw

`$ liquidityd query bank balances cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098`

```
balances:
- amount: "1000000"
  denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
- amount: "799850666"
  denom: stake
- amount: "899849331"
  denom: token
```

## Query

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

Flags:
  -h, --help   help for liquidity

Global Flags:
      --chain-id string    The network chain ID
      --home string        directory for config and data (default "/Users/dongsamb/.liquidityapp")
      --log_level string   The logging level in the format of <module>:<level>,... (default "main:info,state:info,statesync:info,*:error")
      --trace              print out full stack trace on errors

Use "liquidityd query liquidity [command] --help" for more information about a command.
```

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### query batch
`$ liquidityd query liquidity batch --help`
```bash
Query details of a liquidity pool batch
Example:
$ liquidity query liquidity batch 1

Usage:
  liquidityd query liquidity batch [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity batch 1`

```bash
batch:
  batch_index: "3"
  begin_height: "221"
  deposit_msg_index: "2"
  executed: false
  pool_id: "1"
  swap_msg_index: "2"
  withdraw_msg_index: "2"
```

### query batches
`$ liquidityd query liquidity batches --help`
```bash
Query details about all liquidity pools batch on a network.
Example:
$ liquidity query liquidity batches

Usage:
  liquidityd query liquidity batches [flags]
```

`$ liquidityd query liquidity batches`
```bash  

pools_batch:
- batch:
    batch_index: "3"
    begin_height: "221"
    deposit_msg_index: "2"
    executed: false
    pool_id: "1"
    swap_msg_index: "2"
    withdraw_msg_index: "2"
- batch:
    batch_index: "1"
    begin_height: "0"
    deposit_msg_index: "1"
    executed: false
    pool_id: "2"
    swap_msg_index: "1"
    withdraw_msg_index: "1"
pagination:
  next_key: null
  total: "2"
```

### query deposits
`$ liquidityd query liquidity deposits --help`
```bash

Query for all deposit messages on the batch of the liquidity pool specified pool-id

if batch messages are normally processed and from the endblock,
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ liquidity query liquidity deposits 1

Usage:
  liquidityd query liquidity deposits [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity deposits 1 --output json`

```json
{
  "deposits": [
    {
      "msg_height": "51",
      "msg_index": "1",
      "executed": true,
      "succeeded": true,
      "to_be_delete": true,
      "Msg": {
        "depositor_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
        "pool_id": "1",
        "deposit_coins": [
          {
            "denom": "stake",
            "amount": "50000000"
          },
          {
            "denom": "token",
            "amount": "50000000"
          }
        ]
      }
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

empty case

`$ liquidityd query liquidity deposits`

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
$ liquidity query liquidity pool 1

Usage:
  liquidityd query liquidity pool [pool-id] [flags]
```

example query command with result 

`$ liquidityd query liquidity pool 1`
 
```bash
liquidity_pool:
  pool_coin_denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
  pool_id: "1"
  pool_type_id: 1
  reserve_account_address: cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u
  reserve_coin_denoms:
  - stake
  - token
liquidity_pool_batch:
  batch_index: "3"
  begin_height: "221"
  deposit_msg_index: "2"
  executed: false
  pool_id: "1"
  swap_msg_index: "2"
  withdraw_msg_index: "2"
liquidity_pool_metadata:
  pool_coin_total_supply:
    amount: "1000000"
    denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
  pool_id: "1"
  reserve_coins:
  - amount: "100149334"
    denom: stake
  - amount: "100150669"
    denom: token
```


### query pools
`$ liquidityd query liquidity pools  --help`
```   
Query details about all liquidity pools on a network.
Example:
$ liquidity query liquidity pools

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
- liquidity_pool:
    pool_coin_denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
    pool_id: "1"
    pool_type_id: 1
    reserve_account_address: cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u
    reserve_coin_denoms:
    - stake
    - token
  liquidity_pool_batch:
    batch_index: "3"
    begin_height: "221"
    deposit_msg_index: "2"
    executed: false
    pool_id: "1"
    swap_msg_index: "2"
    withdraw_msg_index: "2"
  liquidity_pool_metadata:
    pool_coin_total_supply:
      amount: "1000000"
      denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
    pool_id: "1"
    reserve_coins:
    - amount: "100149334"
      denom: stake
    - amount: "100150669"
      denom: token
- liquidity_pool:
    pool_coin_denom: pool/4718822520A46E7F657C051A7A18A9E8857D2FB47466C9AD81CE2F5F80C61BCC
    pool_id: "2"
    pool_type_id: 1
    reserve_account_address: cosmos1guvgyffq53h87etuq5d85x9fazzh6ta5tq2rjn
    reserve_coin_denoms:
    - atom
    - stake
  liquidity_pool_batch:
    batch_index: "1"
    begin_height: "0"
    deposit_msg_index: "1"
    executed: false
    pool_id: "2"
    swap_msg_index: "1"
    withdraw_msg_index: "1"
  liquidity_pool_metadata:
    pool_coin_total_supply:
      amount: "1000000"
      denom: pool/4718822520A46E7F657C051A7A18A9E8857D2FB47466C9AD81CE2F5F80C61BCC
    pool_id: "2"
    reserve_coins:
    - amount: "100000000"
      denom: atom
    - amount: "100000000"
      denom: stake
```

### query params

example query command with result 

`$ liquidityd query liquidity params`

```bash
init_pool_coin_mint_amount: "1000000"
liquidity_pool_creation_fee:
- amount: "100000000"
  denom: stake
liquidity_pool_types:
- description: ""
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: DefaultPoolType
  pool_type_id: 1
max_order_amount_ratio: "0.100000000000000000"
min_init_deposit_to_pool: "1000000"
swap_fee_rate: "0.003000000000000000"
unit_batch_size: 1
withdraw_fee_rate: "0.003000000000000000"
```


### query swaps
`$ liquidityd query liquidity swaps --help`
```bash
Query for all swap messages on the batch of the liquidity pool

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.

Example:
$ liquidity query liquidity swaps

Usage:
  liquidityd query liquidity swaps [flags]


```

example query command with result

`$ liquidityd query liquidity swaps --output json`
```json
 {
   "swaps": [
     {
       "msg_height": "80",
       "msg_index": "1",
       "executed": true,
       "succeeded": true,
       "to_be_delete": true,
       "order_expiry_height": "80",
       "exchanged_offer_coin": {
         "denom": "token",
         "amount": "1000"
       },
       "remaining_offer_coin": {
         "denom": "token",
         "amount": "0"
       },
       "reserved_offer_coin_fee": {
         "denom": "token",
         "amount": "0"
       },
       "msg": {
         "swap_requester_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
         "pool_id": "1",
         "pool_type_id": 1,
         "swap_type": 1,
         "offer_coin": {
           "denom": "token",
           "amount": "1000"
         },
         "demand_coin_denom": "stake",
         "offer_coin_fee": {
          "denom": "token",
          "amount": "1"
         },
         "order_price": "0.900000000000000000"
       }
     }
   ],
   "pagination": {
     "next_key": null,
     "total": "1"
   }
 }
```

empty case

`$ liquidityd query liquidity swaps`
```bash 
   
pagination:
  next_key: null
  total: "0"
swaps: []

```


### query withdraws
`$ liquidityd query liquidity withdraws --help`
```bash
Query for all withdraws messages on the batch of the liquidity pool specified pool-id

if batch messages are normally processed and from the endblock,
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ liquidity query liquidity withdraws 1

Usage:
  liquidityd query liquidity withdraws [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity withdraws 1 --output json`

```json
{
  "withdraws": [
    {
      "msg_height": "220",
      "msg_index": "1",
      "executed": true,
      "succeeded": true,
      "to_be_delete": true,
      "msg": {
        "withdrawer_address": "cosmos1ta4236u33x0rswerr9rhu2h4ervd67y0dgy098",
        "pool_id": "1",
        "pool_coin": {
          "denom": "pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07",
          "amount": "500000"
        }
      }
    }
  ],
  "pagination": {
    "next_key": null,
    "total": "1"
  }
}
```

empty case
`$ liquidityd query liquidity withdraws`
```bash 
pagination:
  next_key: null
  total: "0"
withdraw_msgs: []

```

## REST/API

You can check local swagger doc page on `YOUR_API_SERVER(ex:127.0.0.1:1317)/swagger-liquidity/` if set `swagger = true` from `app.toml`
or see on [public swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/2.0.2)

According to [migrating-to-new-rest-endpoints](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints), the POST endpoints of the New gGPC-gateway REST are N/A and guided directly to use Protobuf, need to use `cli` or `localhost:1317/cosmos/tx/v1beta1/txs` for broadcast txs temporarily

example of broadcasting txs using the [new REST endpoint (via gRPC-gateway, beta1)](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints)

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```
