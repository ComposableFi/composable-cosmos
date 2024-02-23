echo "********** Running Pre-Scripts **********"

BINARY=$1 
CHAIN_DIR=$(pwd)/mytestnet

KEY="test0"
KEY1="test1"
KEY2="test2"

echo "binary value: $BINARY"
echo "pwd: $(pwd)"
COUNTER_CONTRACT_DIR=$(pwd)/scripts/upgrade/contracts/counter.wasm

KEY_ADDRESS=$($BINARY keys show $KEY -a --keyring-backend test --home $CHAIN_DIR)
echo "key address: $KEY_ADDRESS"

# QUERY balances
$BINARY q bank balances $KEY_ADDRESS --home $CHAIN_DIR

## Deploy the counter contract 
$BINARY tx wasm store $COUNTER_CONTRACT_DIR --from $KEY1 --gas 3000000 --gas-prices 0.025pica --gas-adjustment 1.5 --keyring-backend test --chain-id localpica --home $CHAIN_DIR -y


## interact with the contract. Call the increment function -> counter ++
 


