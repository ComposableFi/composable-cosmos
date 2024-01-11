#!/bin/bash

KEY="mykey"
CHAINID="test-1"
MONIKER="localtestnet"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"
BANKSY_HOME=~/.banksy
# to trace evm
#TRACE="--trace"
TRACE=""

# remove existing daemon
rm -rf ~/.banksy*

layerd config keyring-backend $KEYRING
layerd config chain-id $CHAINID

# if $KEY exists it should be deleted
echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | layerd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover --home $BANKSY_HOME

layerd init $MONIKER --chain-id $CHAINID --home $BANKSY_HOME

# Allocate genesis accounts (cosmos formatted addresses)
layerd add-genesis-account $KEY 100000000000000000000000000stake --keyring-backend $KEYRING --home $BANKSY_HOME

# Sign genesis transaction
layerd gentx $KEY 1000000000000000000000stake --keyring-backend $KEYRING --chain-id $CHAINID --home $BANKSY_HOME

# Collect genesis tx
layerd collect-gentxs --home $BANKSY_HOME

# Run this to ensure everything worked and that the genesis file is setup correctly
layerd validate-genesis --home $BANKSY_HOME

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# update request max size so that we can upload the light client
# '' -e is a must have params on mac, if use linux please delete before run
sed -i'' -e 's/max_body_bytes = /max_body_bytes = 1/g' ~/.banksy/config/config.toml

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
layerd start --pruning=nothing  --minimum-gas-prices=0.0001stake --rpc.laddr tcp://127.0.0.1:36657 --grpc.address localhost:1090 --p2p.laddr tcp://0.0.0.0:36656 --api.address tcp://localhost:2317 --rpc.pprof_laddr tcp://127.0.0.1:7060 --grpc-web.address localhost:1091 --home $BANKSY_HOME
# layerd start --pruning=nothing  --minimum-gas-prices=0.0001stake --rpc.laddr tcp://0.0.0.0:26657 --rpc.grpc_laddr tcp://0.0.0.0:9090 --home $BANKSY_HOME
