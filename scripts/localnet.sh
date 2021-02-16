#!/bin/sh

# Set localnet settings
BINARY=liquidityd
CHAINID=localnet
CHAINDIR=./data
RPCPORT=26657
GRPCPORT=9090
MNEMONIC1="guard cream sadness conduct invite crumble clock pudding hole grit liar hotel maid produce squeeze return argue turtle know drive eight casino maze host"
MNEMONIC2="friend excite rough reopen cover wheel spoon convince island path clean monkey play snow number walnut pull lock shoot hurry dream divide concert discover"

# Stop liquidityd if already running and remove previous data
killall liquidityd
rm -rf $CHAINDIR/$CHAINID

# Add directory for chain, exit if error
if ! mkdir -p $CHAINDIR/$CHAINID 2>/dev/null; then
    echo "Failed to create chain folder. Aborting..."
    exit 1
fi

# Initialize liquidityd with "localnet" chain id
echo "Initializing $CHAINID..."
$BINARY --home $CHAINDIR/$CHAINID init test --chain-id=$CHAINID

echo "Adding genesis accounts..."
echo $MNEMONIC1 | $BINARY --home $CHAINDIR/$CHAINID keys add user1 --recover --keyring-backend=test 
echo $MNEMONIC2 | $BINARY --home $CHAINDIR/$CHAINID keys add validator --recover --keyring-backend=test 
$BINARY --home $CHAINDIR/$CHAINID add-genesis-account $($BINARY keys show validator --keyring-backend test -a) 2000000000stake,1000000000token
$BINARY --home $CHAINDIR/$CHAINID add-genesis-account $($BINARY keys show user1 --keyring-backend test -a) 1000000000stake,1000000000atom

echo "Creating and collecting gentx..."
$BINARY --home $CHAINDIR/$CHAINID gentx validator --amount 1000000000stake --chain-id $CHAINID --keyring-backend test
$BINARY --home $CHAINDIR/$CHAINID collect-gentxs

# Set proper defaults and change ports (MacOS)
echo "Change settings in config.toml file..."
sed -i '' 's#"tcp://127.0.0.1:26657"#"tcp://0.0.0.0:'"$RPCPORT"'"#g' $CHAINDIR/$CHAINID/config/config.toml
sed -i '' 's/timeout_commit = "5s"/timeout_commit = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i '' 's/timeout_propose = "3s"/timeout_propose = "1s"/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i '' 's/index_all_keys = false/index_all_keys = true/g' $CHAINDIR/$CHAINID/config/config.toml
sed -i '' 's/enable = false/enable = true/g' $CHAINDIR/$CHAINID/config/app.toml
sed -i '' 's/swagger = false/swagger = true/g' $CHAINDIR/$CHAINID/config/app.toml

# Start the gaia
echo "Starting $CHAINID in $CHAINDIR..."
echo "Log file is located at $CHAINDIR/$CHAINID.log"
$BINARY --home $CHAINDIR/$CHAINID start --pruning=nothing --grpc.address="0.0.0.0:$GRPCPORT" > $CHAINDIR/$CHAINID.log 2>&1 &