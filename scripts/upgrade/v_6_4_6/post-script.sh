echo "********** Running Post-Scripts **********"

BINARY=$1 
DENOM=${2:-pica}
CHAIN_DIR=$(pwd)/mytestnet

KEY="test0"
KEY1="test1"
KEY2="test2"

WALLET_1=$($BINARY keys show $KEY1 -a --keyring-backend test --home $CHAIN_DIR)  
echo "wallet 1: $WALLET_1"

DEFAULT_GAS_FLAG="--gas 3000000 --gas-prices 0.025$DENOM --gas-adjustment 1.5"
DEFAULT_ENV_FLAG="--keyring-backend test --chain-id localpica --home $CHAIN_DIR"

echo "binary value: $BINARY"
COUNTER_CONTRACT_DIR=$(pwd)/scripts/upgrade/contracts/counter.wasm


## TODO: the old contract address would not work, need to derive a new one 
echo "Old bench32 Contract address: $CONTRACT_ADDRESS"

## Get contract by $CODE_ID
CODE_ID=1 ## TODO: hardfix for now to get the contract, and overide the contract address
CONTRACT_ADDRESS=$($BINARY query wasm list-contract-by-code $CODE_ID -o json | jq -r '.contracts[0]') 
echo "Query contract address: $CONTRACT_ADDRESS"

## Execute the contract, increment counter to 1
$BINARY tx wasm execute $CONTRACT_ADDRESS '{"increment":{}}' --from $KEY1 $DEFAULT_ENV_FLAG $DEFAULT_GAS_FLAG -y -o json > /dev/null

## assert counter value to be 1
sleep 1
COUNTER_VALUE=$($BINARY query wasm contract-state smart $CONTRACT_ADDRESS '{"get_count":{"addr": "'"$WALLET_1"'"}}' -o json | jq -r '.data.count')
echo "COUNTER_VALUE = $COUNTER_VALUE"
if [ "$COUNTER_VALUE" -ne 2 ]; then
    echo "Assertion failed: Expected counter value to be 2, got $COUNTER_VALUE"
    exit 1
fi
echo "Assertion passed: Counter value is 2 as expected"


