
# Liquidityd

Implemented tx cli

- [x]  `create-pool`   Create Liquidity pool with the specified pool-type, deposit-coins
- [x]  `deposit`       Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit-coins
- [x]  `swap`          Swap offer submit to the batch to the Liquidity pool with the specified pool-id with offer-coin, order-price, etc
- [x]  `withdraw`      Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool 

Implemented query cli 

- [x]    `batch`       Query details of a liquidity pool batch of the pool
- [x]    `batches`     Query for all liquidity pools batch
- [x]    `deposit`     Query for the deposit message on the batch of the liquidity pool specified pool-id and msg-index
- [x]    `deposits`    Query for all deposit messages on the batch of the liquidity pool specified pool-id
- [x]    `params`      Query the current liquidity parameters information
- [x]    `pool`        Query details of a liquidity pool
- [x]    `pools`       Query for all liquidity pools
- [x]    `swap`        Query for the swap message on the batch of the liquidity pool specified pool-id and msg-index
- [x]    `swaps`       Query for all swap messages on the batch of the liquidity pool specified pool-id
- [x]    `withdraw`    Query for the withdraw message on the batch of the liquidity pool specified pool-id and msg-index
- [x]    `withdraws`   Query for all withdraw messages on the batch of the liquidity pool specified pool-id

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
  create-pool Create Liquidity pool with the specified pool-type, deposit-coins
  deposit     Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit-coins
  swap        Swap offer to the Liquidity pool with the specified the pool info with offer-coin, order-price
  withdraw    Withdraw submit to the batch from the Liquidity pool with the specified pool-id, pool-coin of the pool
```

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### tx create-pool

`$ liquidityd tx liquidity create-pool --help`

```bash
Create Liquidity pool with the specified pool-type-id, deposit-coins for reserve

Example:
$ liquidity tx liquidity create-pool 1 100000000stake,100000000token --from mykey

Currently, only the default pool-type-id 1 is available on this version
the number of deposit-coins must be two in the pool-type-id 1

{"id":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}

Usage:
  liquidityd tx liquidity create-pool [pool-type-id] [deposit-coins] [flags]
```

example tx command with result 

`$ liquidityd tx liquidity create-pool 1 100000000stake,100000000token --from validator --keyring-backend test --chain-id testing -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgCreatePool",
        "pool_creator_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
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
  "height": "6446",
  "txhash": "B21FB4DF389120DC4D7B5DE0DEE7F4C04CE05E7817B451490B54B59F6AF6364E",
  "codespace": "",
  "code": 0,
  "data": "0A0D0A0B6372656174655F706F6F6C",
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
            },
            {
              "key": "sender",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "sender",
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
              "value": "cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u"
            },
            {
              "key": "amount",
              "value": "100000000stake,100000000token"
            },
            {
              "key": "recipient",
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
            },
            {
              "key": "amount",
              "value": "1000000pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07"
            },
            {
              "key": "recipient",
              "value": "cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl"
            },
            {
              "key": "sender",
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
  "gas_used": "156931",
  "tx": null,
  "timestamp": ""
}
```

already exist case, when duplicated request for same create pool

```json
{
  "height": "6608",
  "txhash": "B4A434F9AA283AFEFAE48DDE5C584175F3AABBBEFB083104045EC29C0A736179",
  "codespace": "liquidity",
  "code": 11,
  "data": "",
  "raw_log": "failed to execute message; message index: 0: the pool already exists",
  "logs": [],
  "info": "",
  "gas_wanted": "200000",
  "gas_used": "49288",
  "tx": null,
  "timestamp": ""
}
```

### tx deposit

`$ liquidityd tx liquidity deposit --help  `

```bash 
liquidityd tx liquidity deposit --help 
Deposit submit to the batch of the Liquidity pool with the specified pool-id, deposit-coins for reserve

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
        "depositor_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
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
  "height": "6659",
  "txhash": "4219875F9C95F0174457E954ED3F6CE9C2EABF7F4DE0125347ADA655AB8D36C2",
  "codespace": "",
  "code": 0,
  "data": "0A110A0F6465706F7369745F746F5F706F6F6C",
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
              "value": "2"
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
  "gas_used": "79179",
  "tx": null,
  "timestamp": ""
}
```


### tx swap

`$ liquidityd tx liquidity swap --help`

```bash  
Swap offer to the Liquidity pool with the specified pool-id, swap-type-id demand-coin-denom
with the coin and the price you're offering and current swap-fee-rate

this requests are stacked in the batch of the liquidity pool, not immediately processed and
processed in the endblock at once with other requests.

Example:
$ liquidity tx liquidity swap 2 1 100000000stake token 0.9 0.003 --from mykey

You should request the same each field as the pool.

Must have sufficient balance half the of the swapFee Rate of the offer coin to reserve offer coin fee.

For explicit calculations, you must enter the params.swap_fee_rate value of the current parameter state.

Currently, only the default pool-type-id, swap-type-id 1 is available on this version
The detailed swap algorithm can be found here.
https://github.com/tendermint/liquidity

Usage:
  liquidityd tx liquidity swap [pool-id] [swap-type-id] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate] [flags]

```

example tx command with result 

`$ liquidityd tx liquidity swap 1 1 1000token stake 0.9 0.003 --from validator --chain-id testing --keyring-backend test -y`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgSwapWithinBatch",
        "swap_requester_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
        "pool_id": "1",
        "swap_type_id": 1,
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
  "height": "6677",
  "txhash": "FCA49D6FC23E1B500C5DED27591CF97AC1B69E635384BF5FCC727A016B31881E",
  "codespace": "",
  "code": 0,
  "data": "0A110A0F6465706F7369745F746F5F706F6F6C",
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
              "value": "2"
            },
            {
              "key": "msg_index",
              "value": "3"
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
  "gas_used": "79179",
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

`$ liquidityd query bank balances cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu`

```
balances:
- amount: "2500000"
  denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
- amount: "650000000"
  denom: stake
- amount: "750000000"
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
        "withdrawer_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
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
  "height": "6705",
  "txhash": "16F0167D3BA6C23943A995BE5803C8AE8E9AEFF83F6352D4A345B4BF384B2A9C",
  "codespace": "",
  "code": 0,
  "data": "0A140A1277697468647261775F66726F6D5F706F6F6C",
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
              "value": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu"
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
              "value": "3"
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
  "gas_used": "67694",
  "tx": null,
  "timestamp": ""
}
```

balances after withdraw

`$ liquidityd query bank balances cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu`

```
balances:
- amount: "2000000"
  denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
- amount: "699850000"
  denom: stake
- amount: "799850000"
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
  batches     Query for all liquidity pools batches
  deposit     Query for the deposit message on the batch of the liquidity pool specified pool-id and msg-index
  deposits    Query for all deposit messages on the batch of the liquidity pool specified pool-id
  params      Query the current liquidity parameters information
  pool        Query details of a liquidity pool
  pools       Query for all liquidity pools
  swap        Query for the swap message on the batch of the liquidity pool specified pool-id and msg-index
  swaps       Query for all swap messages on the batch of the liquidity pool specified pool-id
  withdraw    Query for the withdraw message on the batch of the liquidity pool specified pool-id and msg-index
  withdraws   Query for all withdraw messages on the batch of the liquidity pool specified pool-id
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
  begin_height: "6706"
  deposit_msg_index: "4"
  executed: false
  index: "4"
  pool_id: "1"
  swap_msg_index: "1"
  withdraw_msg_index: "2"
```

### query batches
`$ liquidityd query liquidity batches --help`
```bash
Query details about all liquidity pools batches on a network.
Example:
$ liquidity query liquidity batches

Usage:
  liquidityd query liquidity batches [flags]
```

`$ liquidityd query liquidity batches`
```bash  
batches:
- begin_height: "6706"
  deposit_msg_index: "4"
  executed: false
  index: "4"
  pool_id: "1"
  swap_msg_index: "1"
  withdraw_msg_index: "2"
- begin_height: "0"
  deposit_msg_index: "1"
  executed: false
  index: "1"
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
      "msg_height": "6775",
      "msg_index": "4",
      "executed": true,
      "succeeded": true,
      "to_be_deleted": true,
      "msg": {
        "depositor_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
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
$ liquidity query liquidity pool 1

Usage:
  liquidityd query liquidity pool [pool-id] [flags]
```

example query command with result 

`$ liquidityd query liquidity pool 1`
 
```bash
batch:
  begin_height: "6776"
  deposit_msg_index: "5"
  executed: false
  index: "5"
  swap_msg_index: "1"
  withdraw_msg_index: "2"
id: "1"
metadata:
  pool_coin_total_supply:
    amount: "2499625"
    denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
  reserve_coins:
  - amount: "250150000"
    denom: stake
  - amount: "250150000"
    denom: token
pool_coin_denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
reserve_account_address: cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u
reserve_coin_denoms:
- stake
- token
type_id: 1
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
- batch:
    begin_height: "6776"
    deposit_msg_index: "5"
    executed: false
    index: "5"
    swap_msg_index: "1"
    withdraw_msg_index: "2"
  id: "1"
  metadata:
    pool_coin_total_supply:
      amount: "2499625"
      denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
    reserve_coins:
    - amount: "250150000"
      denom: stake
    - amount: "250150000"
      denom: token
  pool_coin_denom: pool/E4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07
  reserve_account_address: cosmos1unfxz7l7q0s3gmmthgwe3yljk0thhg57ym3p6u
  reserve_coin_denoms:
  - stake
  - token
  type_id: 1
- batch:
    begin_height: "0"
    deposit_msg_index: "1"
    executed: false
    index: "1"
    swap_msg_index: "1"
    withdraw_msg_index: "1"
  id: "2"
  metadata:
    pool_coin_total_supply:
      amount: "1000000"
      denom: pool/4718822520A46E7F657C051A7A18A9E8857D2FB47466C9AD81CE2F5F80C61BCC
    reserve_coins:
    - amount: "100000000"
      denom: atom
    - amount: "100000000"
      denom: stake
  pool_coin_denom: pool/4718822520A46E7F657C051A7A18A9E8857D2FB47466C9AD81CE2F5F80C61BCC
  reserve_account_address: cosmos1guvgyffq53h87etuq5d85x9fazzh6ta5tq2rjn
  reserve_coin_denoms:
  - atom
  - stake
  type_id: 1
```

### query params

example query command with result 

`$ liquidityd query liquidity params`

```bash
init_pool_coin_mint_amount: "1000000"
liquidity_pool_creation_fee:
- amount: "100000000"
  denom: stake
max_order_amount_ratio: "0.100000000000000000"
min_init_deposit_amount: "1000000"
pool_types:
- description: ""
  id: 1
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: DefaultPoolType
max_reserve_coin_amount: "0"
swap_fee_rate: "0.003000000000000000"
unit_batch_size: 1
withdraw_fee_rate: "0.003000000000000000"
```


### query swaps
`$ liquidityd query liquidity swaps --help`
```bash
Query for all swap messages on the batch of the liquidity pool specified pool-id

if batch messages are normally processed and from the endblock,
the resulting state is applied and removed the messages from the beginblock in the next block.
to query for past blocks, you can obtain by specifying the block height through the REST/gRPC API of a node that is not pruned

Example:
$ liquidity query liquidity swaps 1

Usage:
  liquidityd query liquidity swaps [pool-id] [flags]
```

example query command with result

`$ liquidityd query liquidity swaps 2 --output json`
```json
 {
   "swaps": [
     {
       "msg_height": "6829",
       "msg_index": "1",
       "executed": true,
       "succeeded": false,
       "to_be_deleted": false,
       "order_expiry_height": "6829",
       "exchanged_offer_coin": {
         "denom": "stake",
         "amount": "0"
       },
       "remaining_offer_coin": {
         "denom": "stake",
         "amount": "1000"
       },
       "reserved_offer_coin_fee": {
         "denom": "stake",
         "amount": "1"
       },
       "msg": {
         "swap_requester_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
         "pool_id": "2",
         "swap_type_id": 1,
         "offer_coin": {
           "denom": "stake",
           "amount": "1000"
         },
         "demand_coin_denom": "atom",
         "offer_coin_fee": {
           "denom": "stake",
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
      "msg_height": "6848",
      "msg_index": "2",
      "executed": true,
      "succeeded": true,
      "to_be_deleted": true,
      "msg": {
        "withdrawer_address": "cosmos13hwd59j2d4tngxgfpp0v248sxwgenvs232aaqu",
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
`$ liquidityd query liquidity withdraws 2`
```bash 
pagination:
  next_key: null
  total: "0"
withdraws: []

```

## REST/API

You can check local swagger doc page on `YOUR_API_SERVER(ex:127.0.0.1:1317)/swagger-liquidity/` if set `swagger = true` from `app.toml`
or see on [public swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/2.0.2)

According to [migrating-to-new-rest-endpoints](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints), the POST endpoints of the New gGPC-gateway REST are N/A and guided directly to use Protobuf, need to use `cli` or `localhost:1317/cosmos/tx/v1beta1/txs` for broadcast txs temporarily

example of broadcasting txs using the [new REST endpoint (via gRPC-gateway, beta1)](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints)

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```
