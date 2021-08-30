---
title: Liquidityd 
description: A high-level overview of how the command line (CLI) and REST API interfaces work for the liquidity module.
---
# Liquidityd 

This document provides a high-level overview of how the command line (CLI) and REST API interfaces work for the liquidity module.

## Transaction Command Line Interface

- [MsgCreatePool](#msgcreatepool)
  - Create liquidity pool
- [MsgDepositWithinBatch](#msgdepositwithinbatch)
  - Deposit to the liquidity pool batch
- [MsgWithdrawWithinBatch](#msgwithdrawwithinbatch)
  - Withdraw pool coin from the liquidity pool
- [MsgSwapWithinBatch](#msgswapwithinbatch)
  - Swap offer coin with demand coin from the liquidity pool with the given order price

For error codes with the description, see [errors.go](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go).

## Query Command Line Interface

- [Params](#params)
  - Query the current liquidity parameters information
- [Pool](#pool)
  - Query details of a liquidity pool
- [Pools](#pools)
  - Query for all liquidity pools
- [Batch](#batch)
  - Query details of a liquidity pool batch
- [Deposit](#deposit)
  - Query for the deposit message on the batch of the liquidity pool
- [Deposits](#deposits)
  - Query for all deposit messages on the batch of the liquidity pool
- [Withdraw](#withdraw)
  - Query for the withdraw message on the batch of the liquidity pool
- [Withdraws](#withdraws)
  - Query for all withdraw messages on the batch of the liquidity pool
- [Swap](#swap)
  - Query for the swap message on the batch of the liquidity pool
- [Swaps](#swaps)
  - Query for all swap messages on the batch of the liquidity pool

For error codes with the description, see [errors.go](https://github.com/tendermint/liquidity/blob/develop/x/liquidity/types/errors.go).

## REST

A node exposes the REST server default port of `1317`. Configure the port in `[api]` section of the `app.toml` file located in your `$HOME/.liquidityd/config/` directory. When `swagger` param is set to `true`, you can open up your browser and check out the Swagger documentation in `http://localhost:1317/swagger-liquidity/`. You can also reference the public api documentation in [this link](https://app.swaggerhub.com/apis-docs/bharvest/cosmos-sdk_liquidity_module_rest_and_g_rpc_gateway_docs/).

The POST endpoints of the new gGPC-gateway REST are not available. The [Migrating to New REST Endpoints](https://docs.cosmos.network/master/migrations/rest.html#migrating-to-new-rest-endpoints) Cosmos SDK guide suggests to use Protobuf directly. You can use the command line interface or use the temporarily available REST API at `localhost:1317/cosmos/tx/v1beta1/txs`.

For example, to broadcast a transaction by using the [New gRPC-gateway REST Endpoint](https://github.com/cosmos/cosmos-sdk/blob/master/docs/migrations/rest.md#migrating-to-new-rest-endpoints): 

```bash
curl --header "Content-Type: application/json" --request POST --data '{"tx_bytes":"CoMBCoABCh0vdGVuZGVybWludC5saXF1aWRpdHkuTXNnU3dhcBJfCi1jb3Ntb3MxN3dncHpyNGd2YzN1aHBmcnUyNmVhYTJsc203NzJlMnEydjBtZXgQAhgBIAEqDQoFc3Rha2USBDEwMDAyBGF0b206EzExNTAwMDAwMDAwMDAwMDAwMDASWApQCkYKHy9jb3Ntb3MuY3J5cHRvLnNlY3AyNTZrMS5QdWJLZXkSIwohAqzfoAEi0cFg0zqwBuGNvHml4XJNS3EQuVti8/yGH88NEgQKAgh/GAgSBBDAmgwaQGTRN67x2WYF/L5DsRD3ZY1Kt9cVpg3rW+YbXtihxcB6bJWhMxuFr0u9SnGkCuAgOuLH9YU8ROFUo1gGS1RpTz0=","mode":1}' localhost:1317/cosmos/tx/v1beta1/txs
```

## MsgCreatePool

Example `create-pool` tx command:

```bash
liquidityd tx liquidity create-pool 1 1000000000uatom,50000000000uusd --from user1 --keyring-backend test --chain-id testing -b block -o json -y
```

JSON Structure:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.v1beta1.MsgCreatePool",
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

Result

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
          "type": "coin_received",
          "attributes": [
            {
              "key": "receiver",
              "value": "cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl"
            },
            {
              "key": "amount",
              "value": "40000000stake"
            }
          ]
        },
        {
          "type": "coin_spent",
          "attributes": [
            {
              "key": "spender",
              "value": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft"
            },
            {
              "key": "amount",
              "value": "40000000stake"
            }
          ]
        },
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
              "value": "/tendermint.liquidity.v1beta1.MsgCreatePool"
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
              "value": "cosmos1jv65s3grqf6v6jl3dp4t6c9t9rk99cd88lyufl"
            },
            {
              "key": "sender",
              "value": "cosmos1s6cjfm4djg95jkzsfe490yfc9k6wazx6culyft"
            },
            {
              "key": "amount",
              "value": "40000000stake"
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

## MsgDepositWithinBatch

Example `deposit` tx command: 

```bash
liquidityd tx liquidity deposit 1 100000000uatom,5000000000uusd --from validator --keyring-backend test --chain-id testing -y -b block
```

JSON Structure:

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.v1beta1.MsgDepositWithinBatch",
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

Result:

```json
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
          "type": "coin_received",
          "attributes": [
            {
              "key": "receiver",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "amount",
              "value": "100000000uatom,5000000000uusd"
            }
          ]
        },
        {
          "type": "coin_spent",
          "attributes": [
            {
              "key": "spender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "100000000uatom,5000000000uusd"
            }
          ]
        },
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
              "value": "/tendermint.liquidity.v1beta1.MsgDepositWithinBatch"
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

## MsgWithdrawWithinBatch

Example `withdraw` tx command:

```bash
liquidityd tx liquidity withdraw 1 10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295 --from validator --chain-id testing --keyring-backend test -b block -o json -y
```

JSON Structure

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

Result:

```json
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
          "type": "coin_received",
          "attributes": [
            {
              "key": "receiver",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "amount",
              "value": "10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
            }
          ]
        },
        {
          "type": "coin_spent",
          "attributes": [
            {
              "key": "spender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "10000pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295"
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "/tendermint.liquidity.v1beta1.MsgWithdrawWithinBatch"
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
## MsgSwapWithinBatch

Example `swap` tx command:

```bash
liquidityd tx liquidity swap 1 1 50000000uusd uatom 0.019 0.003 --from validator --chain-id testing --keyring-backend test -b block -o json -y
```

JSON Structure:

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

Result:

```json
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
          "type": "coin_received",
          "attributes": [
            {
              "key": "receiver",
              "value": "cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962"
            },
            {
              "key": "amount",
              "value": "50075000uusd"
            }
          ]
        },
        {
          "type": "coin_spent",
          "attributes": [
            {
              "key": "spender",
              "value": "cosmos1h6ht09xx0ue0fqmezk7msgqcc9k20a5x5ynvc3"
            },
            {
              "key": "amount",
              "value": "50075000uusd"
            }
          ]
        },
        {
          "type": "message",
          "attributes": [
            {
              "key": "action",
              "value": "/tendermint.liquidity.v1beta1.MsgSwapWithinBatch"
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
              "value": "50075000uusd"
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
## Params

Example `params` query command:

```bash
$ liquidityd query liquidity params
```

Result:

```json
circuit_breaker_enabled: false
init_pool_coin_mint_amount: "1000000"
max_order_amount_ratio: "0.100000000000000000"
max_reserve_coin_amount: "0"
min_init_deposit_amount: "1000000"
pool_creation_fee:
- amount: "40000000"
  denom: stake
pool_types:
- description: Standard liquidity pool with pool price function X/Y, ESPM constraint,
    and two kinds of reserve coins
  id: 1
  max_reserve_coin_num: 2
  min_reserve_coin_num: 2
  name: StandardLiquidityPool
swap_fee_rate: "0.003000000000000000"
unit_batch_height: 1
withdraw_fee_rate: "0.000000000000000000"
```

## Pool

Example `pool` query command using `pool-id` argument:

```bash
$ liquidityd query liquidity pool 1
```

Result:

```json
pool:
  id: "1"
  pool_coin_denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
  reserve_account_address: cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
  reserve_coin_denoms:
  - uatom
  - uusd
  type_id: 1
```

Example `pool` query command using `--pool-coin-denom` flag:

```bash
$ liquidityd query liquidity pool --pool-coin-denom=pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
```

Result:

```json
pool:
  id: "1"
  pool_coin_denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
  reserve_account_address: cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
  reserve_coin_denoms:
  - uatom
  - uusd
  type_id: 1
```

Example `pool` query command using `--reserve-acc` flag:

```bash
$ liquidityd query liquidity pool --reserve-acc=cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
```

Result:

```json
pool:
  id: "1"
  pool_coin_denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
  reserve_account_address: cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
  reserve_coin_denoms:
  - uatom
  - uusd
  type_id: 1
```

Query reserve coins of the pool:

```bash
$ liquidityd query bank balances cosmos1jmhkafh94jpgakr735r70t32sxq9wzkayzs9we
```

Result:

```json
balances:
- amount: "999003494"
  denom: uatom
- amount: "50050075000"
  denom: uusd
pagination:
  next_key: null
  total: "0"
```

Query total supply of the pool coin

`$ liquidityd query bank total --denom=pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295`

Result:

```json
amount: "1000000"
denom: pool96EF6EA6E5AC828ED87E8D07E7AE2A8180570ADD212117B2DA6F0B75D17A6295
```

## Pools

Example `pools` query command:

```bash
$ liquidityd query liquidity pools
```

Result:

```json
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

## Batch

Example `batch` query command:

```bash
$ liquidityd query liquidity batch 1
```

Result:

```json
batch:
  begin_height: "563"
  deposit_msg_index: "2"
  executed: false
  index: "3"
  pool_id: "1"
  swap_msg_index: "2"
  withdraw_msg_index: "2"
```

## Deposit

Example `deposit` query command:

```bash
$ liquidityd query liquidity deposit 1 1
```

Result:

```json
deposit:
  executed: true
  msg:
    deposit_coins:
    - amount: "1000000000"
      denom: uatom
    - amount: "50000000000"
      denom: uusd
    depositor_address: cosmos1le0a8y0ha99txx0ngsh0qhyyx7cwnjmmju52ed
    pool_id: "1"
  msg_height: "35"
  msg_index: "2"
  succeeded: true
  to_be_deleted: true
```

## Deposits

Example `deposits` query command:

```bash
$ liquidityd query liquidity deposits 1
```

Result:

```json
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

## Withdraw

Example `withdraw` query command:

```bash
$ liquidityd query liquidity withdraws 1 2
```

Result:

```json
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

## Withdraws

Example `withdraws` query command 

```bash
$ liquidityd query liquidity withdraws 1
```

Result:

```json
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
## Swap

Example `swap` query command:

```bash
$ liquidityd query liquidity swaps 1 2
```

Result:

```json
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
## Swaps

Example `swaps` query command:

```bash
$ liquidityd query liquidity swaps 1
```

Result:

```json
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




