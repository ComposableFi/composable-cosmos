echo -e  "\n ********** Running Post-Scripts **********"

BINARY=$1 
DENOM=${2:-upica}
CHAIN_DIR=$(pwd)/mytestnet

KEY="test0"
KEY1="test1"
KEY2="test2"



DEFAULT_GAS_FLAG="--gas 3000000 --gas-prices 0.025$DENOM --gas-adjustment 1.5"
DEFAULT_ENV_FLAG="--keyring-backend test --chain-id localpica --home $CHAIN_DIR"


WALLET_1=$($BINARY keys show $KEY1 -a --keyring-backend test --home $CHAIN_DIR)  
BALANCE_1=$($BINARY query bank balances $WALLET_1 --home $CHAIN_DIR -o json | jq -r '.balances[0].amount')
echo "wallet 1: $WALLET_1 - balance: $BALANCE_1"


WALLET_2=$($BINARY keys show $KEY2 -a --keyring-backend test --home $CHAIN_DIR)
echo "wallet 2: $WALLET_2"



echo "binary value: $BINARY"
COUNTER_CONTRACT_DIR=$(pwd)/scripts/upgrade/contracts/counter.wasm


# ## TODO: the old contract address would not work, need to derive a new one 
# echo "Old bench32 Contract address: $CONTRACT_ADDRESS"

## Get contract by $CODE_ID
echo -e "\n Fetching the new contract address (it got changed after the upgrade)"
CODE_ID=1 ## TODO: hardfix for now to get the contract, and overide the contract address
CONTRACT_ADDRESS=$($BINARY query wasm list-contract-by-code $CODE_ID -o json | jq -r '.contracts[0]') 
echo "Query contract address: $CONTRACT_ADDRESS"

## Fetch code info 
CREATOR=$($BINARY query wasm code-info $CODE_ID -o json | jq -r '.creator')
if [ "$CREATOR" == "$WALLET_1" ]; then
    echo "Assertion passed: Code creator ($CREATOR) is equal to Wallet 1 ($WALLET_1)"
else
    echo "Assertion failed: Code creator ($CREATOR) is not equal to Wallet 1 ($WALLET_1)"
    exit 1
fi


## Fetch contract info 
CONTRACT_CREATOR=$($BINARY query wasm contract $CONTRACT_ADDRESS -o json | jq -r '.contract_info.creator')
if [ "$CONTRACT_CREATOR" == "$WALLET_1" ]; then
    echo "Assertion passed: Contract creator ($CONTRACT_CREATOR) is equal to Wallet 1 ($WALLET_1)"
else
    echo "Assertion failed: Contract creator ($CONTRACT_CREATOR) is not equal to Wallet 1 ($WALLET_1)"
    exit 1
fi

echo -e "\n Testing with new wallet / wallet that has not interacted with the contract before"
## Execute contract with new address
echo "wallet2: init the counter"
$BINARY tx wasm execute $CONTRACT_ADDRESS '{"increment":{}}' --from $KEY2 $DEFAULT_ENV_FLAG $DEFAULT_GAS_FLAG -y -o json > /dev/null # tx1 is to init the counter== 0

sleep 1
echo "wallet2: call the increment() function"
$BINARY tx wasm execute $CONTRACT_ADDRESS '{"increment":{}}' --from $KEY2 $DEFAULT_ENV_FLAG $DEFAULT_GAS_FLAG -y -o json > /dev/null 

sleep 1
echo "wallet2: call the get_count() function"
COUNTER_VALUE_2=$($BINARY query wasm contract-state smart $CONTRACT_ADDRESS '{"get_count":{"addr": "'"$WALLET_2"'"}}' -o json | jq -r '.data.count')
echo "COUNTER_VALUE_2 = $COUNTER_VALUE_2"


echo -e "\n Testing with wallet that has interacted with the contract before"
## Execute the contract, with the existing address. increment counter to 2
echo "wallet1: call the increment() function"
$BINARY tx wasm execute $CONTRACT_ADDRESS '{"increment":{}}' --from $KEY1 $DEFAULT_ENV_FLAG $DEFAULT_GAS_FLAG -y -o json > /dev/null

## assert counter value to be 1
sleep 1
echo "wallet1: call the get_count() function" 
$BINARY query wasm contract-state smart $CONTRACT_ADDRESS '{"get_count":{"addr": "'"$WALLET_1"'"}}' -o json 
COUNTER_VALUE=$($BINARY query wasm contract-state smart $CONTRACT_ADDRESS '{"get_count":{"addr": "'"$WALLET_1"'"}}' -o json | jq -r '.data.count')
echo "COUNTER_VALUE = $COUNTER_VALUE"
# if [ "$COUNTER_VALUE" -ne 2 ]; then
#     echo "Assertion failed: Expected counter value to be 2, got $COUNTER_VALUE"
#     exit 1
# fi
# echo "Assertion passed: Counter value is 2 as expected"


