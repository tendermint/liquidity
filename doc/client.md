
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
- [ ] query with pagination


## Tx

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

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### tx create-pool

`$ ./liquidityd tx liquidity create-pool --help`

```bash
Create Liquidity pool with the specified pool-type-index, deposit coins for reserve

Example:
$ liquidity tx liquidity create-pool 1 100000000acoin,100000000bcoin --from mykey

Currently, only the default pool-type-index 1 is available on this version
the number of deposit coins must be two in the pool-type-index 1

{"pool_type_index":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}

Usage:
  liquidityd tx liquidity create-pool [pool-type-index] [deposit-coins] [flags]


```

example tx command with result 

`$ liquidityd tx liquidity create-pool 1 100000000reservecoin1,100000000reservecoin2 --from node0 --home ./output/node0/liquidityd/ --fees 2stake --chain-id chain-3MYSLc`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgCreateLiquidityPool",
        "pool_creator_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
        "pool_type_index": 1,
        "deposit_coins": [
          {
            "denom": "reservecoin1",
            "amount": "100000000"
          },
          {
            "denom": "reservecoin2",
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
      "amount": [
        {
          "denom": "stake",
          "amount": "2"
        }
      ],
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
  "height": "203",
  "txhash": "BA13B95AEB3AB0FA33E33B64300D59D0C0D846B61242D3F73C9F90AB5B3FFEEA",
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
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "create_liquidity_pool"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "sender",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "module",
              "value": "liquidity"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            // ...
          ]
        },
        {
          "type": "transfer",
          "attributes": [
            {
              "key": "recipient",
              "value": "cosmos1ux8lymc6af2cqzpzshyrjtcurnchlqyqclke67"
            },
            {
              "key": "amount",
              "value": "100000000stake"
            },
            {
              "key": "recipient",
              "value": "cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs"
            },
            {
              "key": "amount",
              "value": "100000000reservecoin1,100000000reservecoin2"
            },
            {
              "key": "recipient",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "amount",
              "value": "1000000cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "160108",
  "tx": null,
  "timestamp": ""
}
```

already exist case, when duplicated request for same create pool

```json
{
  "height": "20",
  "txhash": "2CBA5C6F8C3C3220FA2C5C83C4CDC1314998E4C5632469D6BD7DBF4B16C8C96B",
  "codespace": "liquidity",
  "code": 11,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: the pool already exists",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "56812",
  "tx": null,
  "timestamp": ""
}
```

pool type not exists case

```json
{
  "height": "52",
  "txhash": "7AF58A5C5F416D41976575F354EF79199FC102C19DD3076E02A5DFB8E4A6069E",
  "codespace": "liquidity",
  "code": 2,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: pool type not exists",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "55254",
  "tx": null,
  "timestamp": ""
}
```


### tx deposit

`$ ./liquidityd tx liquidity deposit --help  `

```bash 
./liquidityd tx liquidity deposit --help 
Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit coins for reserve

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity deposit 1 100000000acoin,100000000bcoin --from mykey

You should deposit the same coin as the reserve coin.

Usage:
  liquidityd tx liquidity deposit [pool-id] [deposit-coins] [flags]

```

example tx command with result 


`$ ./liquidityd tx liquidity deposit 1 10000000reservecoin1,10000000reservecoin2 --from node0 --home ./output/node0/liquidityd/ --fees 2stake --chain-id chain-vqZBhx`

```json
 {
   "body": {
     "messages": [
       {
         "@type": "/tendermint.liquidity.MsgDepositToLiquidityPool",
         "depositor_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
         "pool_id": "1",
         "deposit_coins": [
           {
             "denom": "reservecoin1",
             "amount": "10000000"
           },
           {
             "denom": "reservecoin2",
             "amount": "10000000"
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
       "amount": [
         {
           "denom": "stake",
           "amount": "2"
         }
       ],
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
  "height": "1232",
  "txhash": "BE5A788E1BBF5E5DD70D2203AE2E5A1270B2075FAE4ED0DC6842684E9D82B339",
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
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "deposit_to_liquidity_pool"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "module",
              "value": "liquidity"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "batch_id",
              "value": ""
            }
            // ...
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
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "amount",
              "value": "10000000reservecoin1,10000000reservecoin2"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "78915",
  "tx": null,
  "timestamp": ""
}
```


### tx swap

`$ ./liquidityd tx liquidity swap --help`

```bash  
Swap offer to the Liquidity pool with the specified pool-id, swap-type,
demand-coin-denom with the coin and the price you're offering

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity swap 2 1 1 100000000acoin bcoin 1.15 --from mykey

You should request the same each field as the pool.

Currently, only the default swap-type 1 is available on this version
The detailed swap algorithm can be found here.
https://github.com/tendermint/liquidity

```

example tx command with result 

`$ ./liquidityd tx liquidity swap 1 1 1 100000reservecoin1 reservecoin2 1.15 --from node0 --home ./output/node0/liquidityd/ --fees 2stake --chain-id chain-vqZBhx`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgSwap",
        "swap_requester_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
        "pool_id": "1",
        "pool_type_index": 1,
        "swap_type": 1,
        "offer_coin": {
          "denom": "reservecoin1",
          "amount": "100000"
        },
        "demand_coin_denom": "reservecoin2",
        "order_price": "1.150000000000000000"
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
      "amount": [
        {
          "denom": "stake",
          "amount": "2"
        }
      ],
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
  "height": "1863",
  "txhash": "04FC3FD99AC82AAB01CE500B82B4D9A916270C2E18274E9157CA2A07BE70EB5C",
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
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            }
          ]
          //...
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
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "amount",
              "value": "100000reservecoin1"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "75160",
  "tx": null,
  "timestamp": ""
}
```

### tx withdraw

`$ ./liquidityd tx liquidity withdraw --help  `

```bash 
Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool

this requests are stacked in the batch of the liquidity pool, not immediately processed and 
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity withdraw 1 1000cosmos1d9w9j3rq5aunkrkdm86paduz4attl78thlj07f --from mykey

You should request the matched pool-coin as the pool.

Usage:
  liquidityd tx liquidity withdraw [pool-id] [pool-coin] [flags]
```

example tx command with result 

`$ ./liquidityd tx liquidity withdraw 1 1000cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs --from node0 --home ./output/node0/liquidityd/ --fees 2stake --chain-id chain-vqZBhx`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgWithdrawFromLiquidityPool",
        "withdrawer_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
        "pool_id": "1",
        "pool_coin": {
          "denom": "cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs",
          "amount": "1000"
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
      "amount": [
        {
          "denom": "stake",
          "amount": "2"
        }
      ],
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
  "height": "1804",
  "txhash": "C8439CF2C74221DD310069321222C9F7ADFC2E06764A6197E7D2983554BC723C",
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
              "value": "withdraw_from_liquidity_pool"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "module",
              "value": "liquidity"
            },
            {
              "key": "sender",
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            //...
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
              "value": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun"
            },
            {
              "key": "amount",
              "value": "1000cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs"
            }
          ]
        }
      ]
    }
  ],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "72194",
  "tx": null,
  "timestamp": ""
}
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
`$ ./liquidityd query liquidity batch --help`
```bash
Query details of a liquidity pool batch
Example:
$ liquidity query liquidity batch 1

Usage:
  liquidityd query liquidity batch [pool-id] [flags]
```

example query command with result

`$ ./liquidityd query liquidity batch 1`

```bash
liquidity_pool_batch:
  batch_index: "4"
  begin_height: "1864"
  deposit_msg_index: "3"
  executed: false
  pool_id: "1"
  swap_msg_index: "2"
  withdraw_msg_index: "2"
```

### query batches
`$ ./liquidityd query liquidity batches --help`
```bash
Query details about all liquidity pools batch on a network.
Example:
$ liquidity query liquidity batches

Usage:
  liquidityd query liquidity batches [flags]
```

`$ ./liquidityd query liquidity batches`
```bash  
liquidity_pools_batch_response:
- liquidity_pool_batch:
    batch_index: "4"
    begin_height: "1864"
    deposit_msg_index: "3"
    executed: false
    pool_id: "1"
    swap_msg_index: "2"
    withdraw_msg_index: "2"
- liquidity_pool_batch:
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
`$ ./liquidityd query liquidity deposits --help`
```bash

Query for all deposit messages on the batch of the liquidity pool

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.

Example:
$ liquidity query liquidity deposits

```

example query command with result

`$ ./liquidityd query liquidity deposits --output json`
```json
{
  "deposit_msgs": [
    {
      "msg_height": "1232",
      "msg_index": "1",
      "executed": true,
      "succeeded": true,
      "to_be_delete": true,
      "Msg": {
        "depositor_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
        "pool_id": "1",
        "deposit_coins": [
          {
            "denom": "reservecoin1",
            "amount": "10000000"
          },
          {
            "denom": "reservecoin2",
            "amount": "10000000"
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

`$ ./liquidityd query liquidity deposits`
```bash 
deposit_msgs: []
pagination:
  next_key: null
  total: "0"

```

### query pool
`$ ./liquidityd query liquidity pool  --help`
```   
Query details of a liquidity pool
Example:
$ liquidity query liquidity pool 1

Usage:
  liquidityd query liquidity pool [pool-id] [flags]
```

example query command with result 

`./liquidityd query liquidity pool 1`
 
```bash
liquidity_pool:
  pool_coin_denom: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
  pool_id: "1"
  pool_type_index: 1
  reserve_account_address: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
  reserve_coin_denoms:
  - reservecoin1
  - reservecoin2
liquidity_pool_batch:
  batch_index: "4"
  begin_height: "1864"
  deposit_msg_index: "3"
  executed: false
  pool_id: "1"
  swap_msg_index: "2"
  withdraw_msg_index: "2"
liquidity_pool_meta_data:
  pool_coin_total_supply:
    amount: "1199000"
    denom: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
  pool_id: "1"
  reserve_coins:
  - amount: "120000000"
    denom: reservecoin1
  - amount: "119800467"
    denom: reservecoin2
```


### query pools
`$ ./liquidityd query liquidity pools  --help`
```   
./liquidityd query liquidity pools --help   
Query details about all liquidity pools on a network.
Example:
$ liquidity query liquidity pools

Usage:
  liquidityd query liquidity pools [flags]


```

example query command with result 

`./liquidityd query liquidity pools`
 
```bash
./liquidityd query liquidity pools       
liquidity_pools_response:
- liquidity_pool:
    pool_coin_denom: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
    pool_id: "1"
    pool_type_index: 1
    reserve_account_address: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
    reserve_coin_denoms:
    - reservecoin1
    - reservecoin2
  liquidity_pool_batch:
    batch_index: "4"
    begin_height: "1864"
    deposit_msg_index: "3"
    executed: false
    pool_id: "1"
    swap_msg_index: "2"
    withdraw_msg_index: "2"
  liquidity_pool_meta_data:
    pool_coin_total_supply:
      amount: "1199000"
      denom: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
    pool_id: "1"
    reserve_coins:
    - amount: "120000000"
      denom: reservecoin1
    - amount: "119800467"
      denom: reservecoin2
- liquidity_pool:
    pool_coin_denom: cosmos1d9w9j3rq5aunkrkdm86paduz4attl78thlj07f
    pool_id: "2"
    pool_type_index: 1
    reserve_account_address: cosmos1d9w9j3rq5aunkrkdm86paduz4attl78thlj07f
    reserve_coin_denoms:
    - reservecoin1
    - stake
  liquidity_pool_batch:
    batch_index: "1"
    begin_height: "0"
    deposit_msg_index: "1"
    executed: false
    pool_id: "2"
    swap_msg_index: "1"
    withdraw_msg_index: "1"
  liquidity_pool_meta_data:
    pool_coin_total_supply:
      amount: "1000000"
      denom: cosmos1d9w9j3rq5aunkrkdm86paduz4attl78thlj07f
    pool_id: "2"
    reserve_coins:
    - amount: "50000000"
      denom: reservecoin1
    - amount: "1000000"
      denom: stake
pagination:
  next_key: null
  total: "2"
```

### query params

example query command with result 

`./liquidityd query liquidity params`

```bash
./liquidityd query liquidity params
init_pool_coin_mint_amount: "1000000"
liquidity_pool_creation_fee:
- amount: "100000000"
  denom: stake
liquidity_pool_types:
- description: ""
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: DefaultPoolType
  pool_type_index: 1
min_init_deposit_to_pool: "1000000"
swap_fee_rate: "0.003000000000000000"
```


### query swaps
`$ ./liquidityd query liquidity swaps --help`
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

`$ ./liquidityd query liquidity swaps --output json`
```json
 {
   "swap_msgs": [
     {
       "msg_height": "1863",
       "msg_index": "1",
       "executed": true,
       "succeeded": true,
       "to_be_delete": true,
       "order_expiry_height": "1863",
       "exchanged_offer_coin": {
         "denom": "reservecoin1",
         "amount": "100000"
       },
       "remaining_offer_coin": {
         "denom": "reservecoin1",
         "amount": "0"
       },
       "msg": {
         "swap_requester_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
         "pool_id": "1",
         "pool_type_index": 1,
         "swap_type": 1,
         "offer_coin": {
           "denom": "reservecoin1",
           "amount": "100000"
         },
         "demand_coin_denom": "reservecoin2",
         "order_price": "1.150000000000000000"
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

`$ ./liquidityd query liquidity swaps`
```bash 
   
pagination:
  next_key: null
  total: "0"
swap_msgs: []

```


### query withdraws
`$ ./liquidityd query liquidity withdraws --help`
```bash
Query for all withdraws messages on the batch of the liquidity pool

if batch messages are normally processed and from the endblock,  
the resulting state is applied and removed the messages from the beginblock in the next block.

Example:
$ liquidity query liquidity withdraws

Usage:
  liquidityd query liquidity withdraws [flags]


```

example query command with result
`$ ./liquidityd query liquidity withdraws --output json`
```json
{
  "withdraw_msgs": [
    {
      "msg_height": "1804",
      "msg_index": "1",
      "executed": true,
      "succeeded": true,
      "to_be_delete": true,
      "msg": {
        "withdrawer_address": "cosmos1e35y69rhrt7y4yce5l5u73sjnxu0l33wvznyun",
        "pool_id": "1",
        "pool_coin": {
          "denom": "cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs",
          "amount": "1000"
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
`$ ./liquidityd query liquidity withdraws`
```bash 
pagination:
  next_key: null
  total: "0"
withdraw_msgs: []

```

## REST/API

You can check local swagger doc page on `YOUR_API_SERVER(ex:127.0.0.1:1317)/swagger-liquidity/` if set `swagger = true` from `app.toml`
or see on [public swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/2.0.1)

According to [migrating-to-new-rest-endpoints](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints), the POST endpoints of the New gGPC-gateway REST are N/A and guided directly to use Protobuf, need to use `cli` or `localhost:1317/cosmos/tx/v1beta1/txs` for broadcast txs temporarily

example of broadcasting txs using the [new REST endpoint (via gRPC-gateway, beta1)](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints)

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```
