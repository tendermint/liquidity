#!/bin/bash

# Set localnet configuration
# Reference localnet script to see which tokens are given to user account in genesis state
BINARY=liquidityd
CHAIN_ID=localnet
CHAIN_DIR=./data
USER_1=user1 

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

# Create liquidity pool
echo "-> Creating liquidity pool..."
liquidityd tx liquidity create-pool 1 100000000stake,100000000token \
--home $CHAIN_DIR/$CHAIN_ID \
--chain-id $CHAIN_ID \
--from $USER_1 \
--keyring-backend test \
--yes
