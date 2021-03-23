
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
  create-pool Create new liquidity pool with the specified pool type and deposit coins
  deposit     Deposit coins to the specified liquidity pool
  swap        Swap offer coin with demand coin from the specified liquidity pool with the given order price
  withdraw    Withdraw pool coin from the specified liquidity pool
```

See [here](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go) error codes with descriptions

### tx create-pool

`$ liquidityd tx liquidity create-pool --help`

```bash
Create new liquidity pool with the specified pool type and deposit coins.

Example:
$ liquidity tx liquidity create-pool 1 1000000000uatom,50000000000uusd --from mykey

In this example, user requests to create new liquidity pool with 100000000stake and 100000000token.
User must create with a combination of coins that are not already exist in the network.
In this version, pool-type-id 1 is only available, which requires two different coins.

{"id":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}

Usage:
  liquidityd tx liquidity create-pool [pool-type-id] [deposit-coins] [flags]
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
Deposit coins to the specified liquidity pool.

This swap request may not be processed immediately since it will be accumulated in the batch of the liquidity pool.
This will be processed with other requests at once in every end of batch.

Example:
$ liquidity tx liquidity deposit 1 100000000stake,100000000token --from mykey

In this example, user requests to deposit 100000000stake and 100000000token to the specified liquidity pool.
User must deposit the same coin denoms as the reserve coins.

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
Swap offer coin with demand coin from the specified liquidity pool with the given order price.

This swap request may not be processed immediately since it will be accumulated in the batch of the liquidity pool.
This will be processed with other requests at once in every end of batch. 
Note that the order of swap requests is ignored since the universal swap price is calculated within every batch to prevent front running.

The requested swap is executed with a swap price calculated from given swap price function of the pool, the current other swap requests and the current liquidity pool coin reserve status.
Swap orders are executed only when execution swap price is equal or better than submitted order price of the swap order.

Example:
$ liquidity liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from mykey

In this example, we assume there exists a liquidity pool with 1000000000uatom and 50000000000uusd.
User requests to swap 50000000uusd for at least 950000uatom with the order price of 0.019 and swap fee rate of 0.003.
User must have sufficient balance half of the swap-fee-rate of the offer coin to reserve offer coin fee.

The order price is the exchange ratio of X/Y where X is the amount of the first coin and Y is the amount of the second coin when their denoms are sorted alphabetically. 
Increasing order price means to decrease the possibility for your request to be processed and end up buying uatom at cheaper price than the pool price.  

For explicit calculations, you must enter the swap-fee-rate value of the current parameter state.
In this version, swap-type-id 1 is only available. The detailed swap algorithm can be found at https://github.com/tendermint/liquidity

Usage:
  liquidityd tx liquidity swap [pool-id] [swap-type-id] [offer-coin] [demand-coin-denom] [order-price] [swap-fee-rate] [flags]
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

This swap request may not be processed immediately since it will be accumulated in the batch of the liquidity pool.
This will be processed with other requests at once in every end of batch. 

Example:
$ liquidity tx liquidity withdraw 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from mykey

In this example, user requests to withdraw 10000 pool coin from the specified liquidity pool. 
User must request the appropriate pool coin from the specified pool.

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
  batch       Query details of a liquidity pool batch of the pool
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
$ liquidity query liquidity pool 1

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
or see on [public swagger api doc](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/2.0.2)

According to [migrating-to-new-rest-endpoints](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints), the POST endpoints of the New gGPC-gateway REST are N/A and guided directly to use Protobuf, need to use `cli` or `localhost:1317/cosmos/tx/v1beta1/txs` for broadcast txs temporarily

example of broadcasting txs using the [new REST endpoint (via gRPC-gateway, beta1)](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints)

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```
