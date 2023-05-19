#!/bin/bash

KEY="mykey"
CHAINID="banksy-testnet-1"
KEYALGO="secp256k1"
KEYRING="test"

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf ~/.banksy*

./banksyd config keyring-backend $KEYRING
./banksyd config chain-id $CHAINID

# if $KEY exists it should be deleted
echo "taste shoot adapt slow truly grape gift need suggest midnight burger horn whisper hat vast aspect exit scorpion jewel axis great area awful blind" | ./banksyd keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover


./banksyd tx 08-wasm push-wasm contracts/ics10_grandpa_cw.wasm --from mykey --keyring-backend test --gas 902152622 --fees 920166stake -y