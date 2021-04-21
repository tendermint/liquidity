#!/bin/sh

# Set localnet settings
BINARY=liquidityd
CHAIN_ID=localnet
CHAIN_DIR=./data
RPC_PORT=26657
GRPC_PORT=9090
MNEMONIC_1="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"
MNEMONIC_2="friend excite rough reopen cover wheel spoon convince island path clean monkey play snow number walnut pull lock shoot hurry dream divide concert discover"
MNEMONIC_3="fuel obscure melt april direct second usual hair leave hobby beef bacon solid drum used law mercy worry fat super must ritual bring faculty"
GENESIS_COINS=10000000000stake,10000000000token,10000000000atom

# Stop liquidityd if it is already running 
if pgrep -x "$BINARY" >/dev/null; then
    echo "Terminating $BINARY..."
    killall liquidityd
fi

# Remove previous data
rm -rf $CHAIN_DIR/$CHAIN_ID

# Add directory for chain, exit if error
if ! mkdir -p $CHAIN_DIR/$CHAIN_ID 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

# Initialize liquidityd with "localnet" chain id
echo "Initializing $CHAIN_ID..."
$BINARY --home $CHAIN_DIR/$CHAIN_ID init test --chain-id=$CHAIN_ID

echo "Adding genesis accounts..."
echo $MNEMONIC_1 | $BINARY --home $CHAIN_DIR/$CHAIN_ID keys add validator --recover --keyring-backend=test 
echo $MNEMONIC_2 | $BINARY --home $CHAIN_DIR/$CHAIN_ID keys add user1 --recover --keyring-backend=test 
echo $MNEMONIC_3 | $BINARY --home $CHAIN_DIR/$CHAIN_ID keys add user2 --recover --keyring-backend=test 
$BINARY --home $CHAIN_DIR/$CHAIN_ID add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAIN_ID keys show validator --keyring-backend test -a) $GENESIS_COINS
$BINARY --home $CHAIN_DIR/$CHAIN_ID add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAIN_ID keys show user1 --keyring-backend test -a) $GENESIS_COINS
$BINARY --home $CHAIN_DIR/$CHAIN_ID add-genesis-account $($BINARY --home $CHAIN_DIR/$CHAIN_ID keys show user2 --keyring-backend test -a) $GENESIS_COINS

echo "Creating and collecting gentx..."
$BINARY --home $CHAIN_DIR/$CHAIN_ID gentx validator 1000000000stake --chain-id $CHAIN_ID --keyring-backend test
$BINARY --home $CHAIN_DIR/$CHAIN_ID collect-gentxs

# Set proper defaults and change ports (MacOS)
echo "Change settings in config.toml file..."
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPC_PORT"'"#g' $CHAIN_DIR/$CHAIN_ID/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAIN_DIR/$CHAIN_ID/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAIN_DIR/$CHAIN_ID/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $CHAIN_DIR/$CHAIN_ID/config/config.toml
sed -i '' 's/enable = false/enable = true/g' $CHAIN_DIR/$CHAIN_ID/config/app.toml
sed -i '' 's/swagger = false/swagger = true/g' $CHAIN_DIR/$CHAIN_ID/config/app.toml

# Start the gaia
echo "Starting $CHAIN_ID in $CHAIN_DIR..."
echo "Log file is located at $CHAIN_DIR/$CHAIN_ID.log"
$BINARY --home $CHAIN_DIR/$CHAIN_ID start --pruning=nothing --grpc.address="0.0.0.0:$GRPC_PORT" > $CHAIN_DIR/$CHAIN_ID.log 2>&1 &