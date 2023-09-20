#!/bin/bash
set -e

BINARY=${BINARY:-centaurid}
BASE_DIR=./data
CHAINID=${CHAINID:-test-1}
STAKEDENOM=${STAKEDENOM:-stake}
CONTRACTS_BINARIES_DIR=${CONTRACTS_BINARIES_DIR:-./contracts}
THIRD_PARTY_CONTRACTS_DIR=${THIRD_PARTY_CONTRACTS_DIR:-./contracts_thirdparty}

# IMPORTANT! minimum_gas_prices should always contain at least one record, otherwise the chain will not start or halt
# ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2 denom is required by intgration tests (test:tokenomics)
MIN_GAS_PRICES_DEFAULT='[{"denom":"ibc/27394FB092D2ECCD56123C74F36E4C1F926001CEADA9CA97EA622B25F41E5EB2","amount":"0"},{"denom":"stake","amount":"0"}]'
MIN_GAS_PRICES=${MIN_GAS_PRICES:-"$MIN_GAS_PRICES_DEFAULT"}


CHAIN_DIR="$BASE_DIR/$CHAINID"
GENESIS_PATH="$CHAIN_DIR/config/genesis.json"

ADMIN_ADDRESS=$($BINARY keys show demowallet1 -a --home "$CHAIN_DIR" --keyring-backend test)
SECOND_MULTISIG_ADDRESS=$($BINARY keys show demowallet2 -a --home "$CHAIN_DIR" --keyring-backend test)

echo "Add consumer section..."
$BINARY add-consumer-section --home "$CHAIN_DIR"
### PARAMETERS SECTION

## slashing params
SLASHING_SIGNED_BLOCKS_WINDOW=140000
SLASHING_MIN_SIGNED=0.050000000000000000
SLASHING_FRACTION_DOUBLE_SIGN=0.010000000000000000
SLASHING_FRACTION_DOWNTIME=0.000100000000000000

function set_genesis_param() {
  param_name=$1
  param_value=$2
  sed -i -e "s;\"$param_name\":.*;\"$param_name\": $param_value;g" "$GENESIS_PATH"
}

set_genesis_param signed_blocks_window        "\"$SLASHING_SIGNED_BLOCKS_WINDOW\","         # slashing
set_genesis_param min_signed_per_window       "\"$SLASHING_MIN_SIGNED\","                   # slashing
set_genesis_param slash_fraction_double_sign  "\"$SLASHING_FRACTION_DOUBLE_SIGN\","         # slashing
set_genesis_param slash_fraction_downtime     "\"$SLASHING_FRACTION_DOWNTIME\""             # slashing
set_genesis_param minimum_gas_prices          "$MIN_GAS_PRICES"                             # globalfee

if ! jq -e . "$GENESIS_PATH" >/dev/null 2>&1; then
    echo "genesis appears to become incorrect json" >&2
    exit 1
fi

echo "$BINARY"
