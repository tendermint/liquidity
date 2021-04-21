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
# liquidityd tx liquidity withdraw 1 1000poolE4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07 --home ./data/localnet --chain-id localnet --from user1 --keyring-backend test --yes
echo "-> Withdrawing coins from the liquidity pool 1..."
$BINARY tx liquidity withdraw 1 1000poolE4D2617BFE03E1146F6BBA1D9893F2B3D77BA29E7ED532BB721A39FF1ECC1B07 \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from user1 \
--keyring-backend test \
--yes

sleep 2

# liquidityd q liquidity withdraw 1 --home ./data/localnet --output json | jq
echo "-> Querying liquidity withdrawals..."
$BINARY q liquidity withdraws 1 \
--home $CHAIN_DIR/$CHAIN_ID \
--output json | jq

# Check the withdraw_msg_index update
# liquidityd q liquidity batch 1 --home ./data/localnet --chain-id localnet --output json | jq
echo "-> Querying details of liquidity pool 1 batch..."
$BINARY q liquidity batch 1  \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--output json | jq

# Ensure the existence of the liquidity pool. 
# If there is no liquidity pool created then use create-pool script to create liquidity pool.
# liquidityd tx liquidity withdraw 1 500pool4718822520A46E7F657C051A7A18A9E8857D2FB47466C9AD81CE2F5F80C61BCC --home ./data/localnet --chain-id localnet --from user1 --keyring-backend test --yes
echo "-> Withdrawing coins from the liquidity pool 2..."
$BINARY tx liquidity withdraw 2 500pool4718822520A46E7F657C051A7A18A9E8857D2FB47466C9AD81CE2F5F80C61BCC \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from user2 \
--keyring-backend test \
--yes

sleep 2

# liquidityd q liquidity withdraws 1 --home ./data/localnet --output json | jq
echo "-> Querying liquidity withdrawals..."
$BINARY q liquidity pools \
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
