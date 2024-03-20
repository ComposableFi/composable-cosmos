#!/bin/bash

KEY="mykey"
CHAINID="test-1"
KEYALGO="secp256k1"
KEYRING="test"

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf ~/.centauri*

~/go/bin/picad config keyring-backend $KEYRING
~/go/bin/picad config chain-id $CHAINID

# if $KEY exists it should be deleted
echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | ~/go/bin/picad  keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover


~/go/bin/picad  tx 08-wasm push-wasm contracts/ics10_grandpa_cw.wasm --from mykey --keyring-backend test --gas 902152622 --fees 920166stake -y