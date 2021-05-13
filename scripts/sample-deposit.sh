#!/bin/bash

# Set localnet configuration
# Reference localnet script to see which tokens are given to user account in genesis state
BINARY=liquidityd
CHAIN_ID=localnet
CHAIN_DIR=./data
USER_1_ADDRESS=cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu
USER_2_ADDRESS=cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny

# Ensure jq is installed
if [[ ! -x "$(which jq)" ]]; then
  echo "jq (a tool for parsing json in the command line) is required..."
  echo "https://stedolan.github.io/jq/download/"
  exit 1
fi

# Ensure liquidityd is installed
if ! [ -x "$(which $BINARY)" ]; then
  echo "Error: liquidityd is not installed. Try building $BINARY by 'make install'" >&2
  exit 1
fi

# Ensure localnet is running
if [[ "$(pgrep $BINARY)" == "" ]];then
    echo "Error: localnet is not running. Try running localnet by 'make localnet" 
    exit 1
fi

# liquidityd q bank balances cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu --home ./data/localnet --output json | jq
echo "-> Checking user1 account balances..."
$BINARY q bank balances $USER_1_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# liquidityd q bank balances cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny --home ./data/localnet --output json | jq
echo "-> Checking user2 account balances..."
$BINARY q bank balances $USER_2_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# liquidityd q liquidity batch 1 --home ./data/localnet --chain-id localnet --output json | jq
echo "-> Querying details of liquidity pool 1 batch..."
$BINARY q liquidity batch 1  \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--output json | jq

# Ensure the existence of the liquidity pool. 
# If there is no liquidity pool created then use create-pool script to create liquidity pool.
# liquidityd tx liquidity deposit 1 100000stake,200000token --home ./data/localnet --chain-id localnet --from user1 --keyring-backend test --yes
echo "-> Depositing coins to the liquidity pool 1..."
$BINARY tx liquidity deposit 1 100000stake,200000token \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from user1 \
--keyring-backend test \
--yes

sleep 2

# liquidityd q liquidity deposits 1 --home ./data/localnet --output json | jq
echo "-> Querying liquidity deposits..."
$BINARY q liquidity deposits 1 \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# Check the deposit_msg_index update
# liquidityd q liquidity batch 1 --home ./data/localnet --chain-id localnet --output json | jq
echo "-> Querying details of liquidity pool 1 batch..."
$BINARY q liquidity batch 1  \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--output json | jq

# Ensure the existence of the liquidity pool. 
# If there is no liquidity pool created then use create-pool script to create liquidity pool.
# liquidityd tx liquidity deposit 2 100000stake,200000token --home ./data/localnet --chain-id localnet --from user2 --keyring-backend test --yes
echo "-> Depositing coins to the liquidity pool 2..."
$BINARY tx liquidity deposit 2 100000stake,200000atom \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from user2 \
--keyring-backend test \
--yes

sleep 2

# liquidityd q liquidity deposits 2 --home ./data/localnet --output json | jq
echo "-> Querying liquidity deposits..."
$BINARY q liquidity deposits 2 \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# liquidityd q bank balances cosmos1mzgucqnfr2l8cj5apvdpllhzt4zeuh2cshz5xu --home ./data/localnet --output json | jq
echo "-> Checking user1 account balances after..."
$BINARY q bank balances $USER_1_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# liquidityd q bank balances cosmos185fflsvwrz0cx46w6qada7mdy92m6kx4gqx0ny --home ./data/localnet --output json | jq
echo "-> Checking user2 account balances after..."
$BINARY q bank balances $USER_2_ADDRESS \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

