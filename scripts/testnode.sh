#!/bin/bash

KEY="mykey"
CHAINID="test-1"
MONIKER="localtestnet"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# remove existing daemon
rm -rf ~/.banksy*

picad config keyring-backend $KEYRING
picad config chain-id $CHAINID

# if $KEY exists it should be deleted
echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | picad keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover

picad init $MONIKER --chain-id $CHAINID

# Allocate genesis accounts (cosmos formatted addresses)
picad add-genesis-account $KEY 100000000000000000000000000stake --keyring-backend $KEYRING

# Sign genesis transaction
picad gentx $KEY 1000000000000000000000stake --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
picad collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
picad validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# update request max size so that we can upload the light client
# '' -e is a must have params on mac, if use linux please delete before run
sed -i'' -e 's/max_body_bytes = /max_body_bytes = 1/g' ~/.banksy/config/config.toml

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
picad start --pruning=nothing  --minimum-gas-prices=0.0001stake --rpc.laddr tcp://0.0.0.0:26657
