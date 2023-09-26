#!/bin/bash

KEY="mykey"
CHAINID="centaurid-t1"
KEYALGO="secp256k1"
KEYRING="test"

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }


centaurid config keyring-backend $KEYRING
centaurid config chain-id $CHAINID

centaurid  tx 08-wasm push-wasm contracts/ics10_grandpa_cw.wasm --from test --keyring-backend test --gas 902152622 --fees 920166stake -y