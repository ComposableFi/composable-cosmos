KEY="mykey"
CHAINID="test-1"
MONIKER="localtestnet"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf ~/.polytope*

polytoped config keyring-backend $KEYRING
polytoped config chain-id $CHAINID

# if $KEY exists it should be deleted
echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | polytoped keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover

polytoped init $MONIKER --chain-id $CHAINID 

# Allocate genesis accounts (cosmos formatted addresses)
polytoped add-genesis-account $KEY 100000000000000000000000000stake --keyring-backend $KEYRING

# Sign genesis transaction
polytoped gentx $KEY 1000000000000000000000stake --keyring-backend $KEYRING --chain-id $CHAINID

# Collect genesis tx
polytoped collect-gentxs

# Run this to ensure everything worked and that the genesis file is setup correctly
polytoped validate-genesis

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
polytoped start --pruning=nothing  --minimum-gas-prices=0.0001stake 
