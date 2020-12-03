
## Liquidityd

```bash
Usage:
  liquidityd [command]

Available Commands:


  add-genesis-account Add a genesis account to genesis.json
  collect-gentxs      Collect genesis txs and output a genesis.json file
  debug               Tool for helping with debugging your application
  export              Export state to JSON
  gentx               Generate a genesis tx carrying a self delegation
  help                Help about any command
  init                Initialize private validator, p2p, genesis, and application configuration files
  keys                Manage your application's keys
  migrate             Migrate genesis to a specified target version
  query               Querying subcommands
  start               Run the full node
  status              Query remote node for status
  tendermint          Tendermint subcommands
  testnet             Initialize files for a liquidityapp testnet
  tx                  Transactions subcommands
  unsafe-reset-all    Resets the blockchain database, removes address book files, and resets priv_validator.json to the genesis state
  validate-genesis    validates the genesis file at the default location or at the location passed as an arg
  version             Print the application binary version information

Flags:
  -h, --help          help for liquidityd
      --home string   directory for config and data (default "/Users/dongsamb/.liquidityapp")
      --trace         print out full stack trace on errors
```


### Tx

`$ liquidityd tx liquidity --help`

```bash
Liquidity transaction subcommands

Usage:
  liquidityd tx liquidity [flags]


Available Commands:
  create-pool Create Liquidity pool with the specified pool-type, deposit coins

  *WIP, More will soon be added.*

Flags:
  -h, --help   help for liquidity

Global Flags:
      --chain-id string   The network chain ID
      --home string       directory for config and data (default "/Users/dongsamb/.liquidityapp")
      --trace             print out full stack trace on errors
```

`$ liquidityd tx liquidity create-pool --help`
```
Create Liquidity pool with the specified pool-type-index, deposit coins for reserve

Example:
$ liquidity tx liquidity create-pool 1 100000000acoin,100000000bcoin --from mykey

Currently, only the default pool-type-index 1 is available
the number of deposit coins must be two in the pool-type-index 1

{"pool_type_index":1,"name":"ConstantProductLiquidityPool","min_reserve_coin_num":2,"max_reserve_coin_num":2,"description":""}

Usage:
  liquidityd tx liquidity create-pool [pool-type-index] [deposit-coins] [flags]
```

### create-pool

example tx command with result 

`$ liquidityd tx liquidity create-pool 1 100000000reservecoin1,100000000reservecoin2 --from node0 --home ./output/node0/liquidityd/ --fees 2stake --chain-id chain-3MYSLc`

```json
{
  "body": {
    "messages": [
      {
        "@type": "/tendermint.liquidity.MsgCreateLiquidityPool",
        "pool_creator_address": "cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l",
        "pool_type_index": 1,
        "reserve_coin_denoms": [
          "reservecoin1",
          "reservecoin2"
        ],
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
```
{
  "height": "46",
  "txhash": "4F8D1EDD3CD3F57A6741FA4C706CF82B7107721002CEF1619EDF97F1C4BA6120",
  "codespace": "",
  "code": 0,
  "data": "0A170A156372656174655F6C69717569646974795F706F6F6C",
  "raw_log": "[{\"events\":[{\"type\":\"message\",\"attributes\":[{\"key\":\"action\",\"value\":\"create_liquidity_pool\"},{\"key\":\"sender\",\"value\":\"cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l\"},{\"key\":\"sender\",\"value\":\"cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l\"},{\"key\":\"sender\",\"value\":\"cosmos1tx68a8k9yz54z06qfve9l2zxvgsz4ka3hr8962\"},{\"key\":\"module\",\"value\":\"liquidity\"},{\"key\":\"sender\",\"value\":\"cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l\"},{\"key\":\"liquidity_pool_id\"},{\"key\":\"liquidity_pool_type_index\",\"value\":\"1\"},{\"key\":\"reserve_coin_denoms\"},{\"key\":\"reserve_account\"},{\"key\":\"pool_coin_denom\"},{\"key\":\"swap_fee_rate\"},{\"key\":\"liquidity_pool_fee_rate\"},{\"key\":\"batch_size\"}]},{\"type\":\"transfer\",\"attributes\":[{\"key\":\"recipient\",\"value\":\"cosmos1ux8lymc6af2cqzpzshyrjtcurnchlqyqclke67\"},{\"key\":\"amount\",\"value\":\"100000000stake\"},{\"key\":\"recipient\",\"value\":\"cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs\"},{\"key\":\"amount\",\"value\":\"100000000reservecoin1,100000000reservecoin2\"},{\"key\":\"recipient\",\"value\":\"cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l\"},{\"key\":\"amount\",\"value\":\"1000000cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs\"}]}]}]",
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
              "value": "cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l"
            },
            {
              "key": "sender",
              "value": "cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l"
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
              "value": "cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l"
            },
            {
              "key": "liquidity_pool_id",
              "value": ""
            },
            {
              "key": "liquidity_pool_type_index",
              "value": "1"
            },
            {
              "key": "reserve_coin_denoms",
              "value": ""
            },
            {
              "key": "reserve_account",
              "value": ""
            },
            {
              "key": "pool_coin_denom",
              "value": ""
            },
            {
              "key": "swap_fee_rate",
              "value": ""
            },
            {
              "key": "liquidity_pool_fee_rate",
              "value": ""
            },
            {
              "key": "batch_size",
              "value": ""
            }
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
              "value": "cosmos1sp5rl2dvwthzm4sujucfwk97k4678de4k0yu9l"
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

### Query

`$ liquidityd query liquidity --help`
```bash
Querying commands for the liquidity module

Usage:
  liquidityd query liquidity [flags]
  liquidityd query liquidity [command]

Available Commands:
  pool           Query details of a liquidity pool
  params         Query the current liquidity parameters information
  *WIP, More will soon be added.*

Flags:
  -h, --help   help for liquidity

Global Flags:
      --chain-id string   The network chain ID
      --home string       directory for config and data (default "/Users/dongsamb/.liquidityapp")
      --trace             print out full stack trace on errors

Use "liquidityd query liquidity [command] --help" for more information about a command.
```

`$ ./liquidityd query liquidity pool  --help   `
```   
Query details of a liquidity pool
Example:
$ liquidity query liquidity pool 1

Usage:
  liquidityd query liquidity pool [pool-id] [flags]

Flags:
      --height int      Use a specific height to query state at (this can error if the node is pruning state)
  -h, --help            help for pool
      --node string     <host>:<port> to Tendermint RPC interface for this chain (default "tcp://localhost:26657")
  -o, --output string   Output format (text|json) (default "text")

Global Flags:
      --chain-id string    The network chain ID
      --home string        directory for config and data (default "/Users/dongsamb/.liquidityapp")
      --log_level string   The logging level in the format of <module>:<level>,... (default "main:info,state:info,statesync:info,*:error")
      --trace              print out full stack trace on errors

```
### pool

example query command with result 

`./liquidityd query liquidity pool 1`
 
```bash
pool_coin_denom: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
pool_id: "1"
pool_type_index: 1
reserve_account_address: cosmos1qz38nymksetqd2d4qesrxpffzywuel82a4l0vs
reserve_coin_denoms:
- reservecoin1
- reservecoin2
```

### params

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


## Export, Genesis State

### export empty state case
`./liquidityd testnet --v 1` 

`./liquidityd start --home ./output/node0/liquidityd/`

`./liquidityd export  --home ./output/node0/liquidityd/`

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

`./liquidityd testnet --v 1`

`./liquidityd start --home ./output/node0/liquidityd/`

`cat output/node0/liquidityd/config/genesis.json | grep chain_id`

`./liquidityd tx liquidity create-pool 1 100000000reservecoin1,100000000reservecoin2 --from node0  --home ./output/node0/liquidityd/ --fees 2stake --chain-id <CHAIN-ID>`

`./liquidityd export --home ./output/node0/liquidityd/`

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