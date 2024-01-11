#!/bin/bash

KEY="mykey"
CHAINID="banksy-testnet-4"
MONIKER="testnet"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

HOMEDIR="$HOME/.centauri-tmp"

CONFIG=$HOMEDIR/config/config.toml
GENESIS=$HOMEDIR/config/genesis.json
TMP_GENESIS=$HOMEDIR/config/tmp_genesis.json

# remove existing daemon
rm -rf $HOMEDIR

centaurid config keyring-backend $KEYRING --home $HOMEDIR
centaurid config chain-id $CHAINID --home $HOMEDIR

# if $KEY exists it should be deleted
centaurid keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --home $HOMEDIR

centaurid init $MONIKER --chain-id $CHAINID --home $HOMEDIR

# Change parameter token denominations to ppica
jq '.app_state["staking"]["params"]["bond_denom"]="ppica"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["staking"]["params"]["unbonding_time"]="300s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["mint"]["params"]["mint_denom"]="ppica"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["mint"]["incentives_supply"]["denom"]="ppica"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["crisis"]["constant_fee"]["denom"]="ppica"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["params"]["min_deposit"][0]["denom"]="ppica"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["params"]["max_deposit_period"]="300s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"
jq '.app_state["gov"]["params"]["voting_period"]="1800s"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

# Set gas limit in genesis
jq '.consensus_params["block"]["max_gas"]="10000000"' "$GENESIS" >"$TMP_GENESIS" && mv "$TMP_GENESIS" "$GENESIS"

if [[ $1 == "pending" ]]; then
    if [[ "$OSTYPE" == "darwin"* ]]; then
        sed -i '' 's/timeout_propose = "3s"/timeout_propose = "30s"/g' "$CONFIG"
        sed -i '' 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' "$CONFIG"
        sed -i '' 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' "$CONFIG"
        sed -i '' 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' "$CONFIG"
        sed -i '' 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' "$CONFIG"
        sed -i '' 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' "$CONFIG"
        sed -i '' 's/timeout_commit = "5s"/timeout_commit = "150s"/g' "$CONFIG"
        sed -i '' 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' "$CONFIG"
    else
        sed -i 's/timeout_propose = "3s"/timeout_propose = "30s"/g' "$CONFIG"
        sed -i 's/timeout_propose_delta = "500ms"/timeout_propose_delta = "5s"/g' "$CONFIG"
        sed -i 's/timeout_prevote = "1s"/timeout_prevote = "10s"/g' "$CONFIG"
        sed -i 's/timeout_prevote_delta = "500ms"/timeout_prevote_delta = "5s"/g' "$CONFIG"
        sed -i 's/timeout_precommit = "1s"/timeout_precommit = "10s"/g' "$CONFIG"
        sed -i 's/timeout_precommit_delta = "500ms"/timeout_precommit_delta = "5s"/g' "$CONFIG"
        sed -i 's/timeout_commit = "5s"/timeout_commit = "150s"/g' "$CONFIG"
        sed -i 's/timeout_broadcast_tx_commit = "10s"/timeout_broadcast_tx_commit = "150s"/g' "$CONFIG"
    fi
fi

if [[ "$OSTYPE" == "darwin"* ]]; then
    sed -i '' 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' "$CONFIG"
else
    sed -i 's/laddr = "tcp:\/\/127.0.0.1:26657"/laddr = "tcp:\/\/0.0.0.0:26657"/g' "$CONFIG"
fi

# Allocate genesis accounts (cosmos formatted addresses)
centaurid add-genesis-account $KEY 100000000000000000000000000ppica --keyring-backend $KEYRING --home $HOMEDIR

# Sign genesis transaction
centaurid gentx $KEY 1000000000000000000000ppica --keyring-backend $KEYRING --chain-id $CHAINID --home $HOMEDIR

# Collect genesis tx
centaurid collect-gentxs --home $HOMEDIR

# Run this to ensure everything worked and that the genesis file is setup correctly
centaurid validate-genesis --home $HOMEDIR

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
centaurid start --pruning=nothing  --minimum-gas-prices=0.0001ppica --rpc.laddr tcp://0.0.0.0:26657 --home $HOMEDIR
