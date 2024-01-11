KEY="mykey"
CHAINID="centauri-testnet-1"
MONIKER="localtestnet"
KEYALGO="secp256k1"
KEYRING="test"
LOGLEVEL="info"
# to trace evm
#TRACE="--trace"
TRACE=""
BANKSY_HOME=~/.banksy

# validate dependencies are installed
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

# remove existing daemon
rm -rf ~/.banksy
alias centaurid=centaurid-nogov-2/centaurid

centaurid config keyring-backend $KEYRING
centaurid config chain-id $CHAINID

# if $KEY exists it should be deleted
#echo "taste shoot adapt slow truly grape gift need suggest midnight burger horn whisper hat vast aspect exit scorpion jewel axis great area awful blind" | centaurid keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover
#echo "sense state fringe stool behind explain area quit ugly affair develop thumb clinic weasel choice atom gesture spare sea renew penalty second upon peace" | centaurid keys add k1 --keyring-backend $KEYRING --algo $KEYALGO --recover
echo "sense state fringe stool behind explain area quit ugly affair develop thumb clinic weasel choice atom gesture spare sea renew penalty second upon peace" | centaurid keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover --home $BANKSY_HOME

centaurid init $MONIKER --chain-id $CHAINID --home $BANKSY_HOME

# Allocate genesis accounts (centauri formatted addresses)
centaurid add-genesis-account $KEY 10000000000000000000stake --keyring-backend $KEYRING --home $BANKSY_HOME
#centaurid add-genesis-account k1 10000000000000000000stake --keyring-backend $KEYRING

# Sign genesis transaction centauri1594tdya20hxz7kjenkn5w09jljyvdfk8kx5rd6
centaurid gentx $KEY 100000000000000000stake --keyring-backend $KEYRING --chain-id $CHAINID --home $BANKSY_HOME

# Collect genesis tx
centaurid collect-gentxs --home $BANKSY_HOME

# Run this to ensure everything worked and that the genesis file is setup correctly
centaurid validate-genesis --home $BANKSY_HOME

if [[ $1 == "pending" ]]; then
  echo "pending mode is on, please wait for the first block committed."
fi

# update request max size so that we can upload the light client
# '' -e is a must have params on mac, if use linux please delete before run
sed -i'' -e 's/max_body_bytes = /max_body_bytes = 1000/g' ~/.banksy/config/config.toml
sed -i'' -e 's/max_tx_bytes = /max_tx_bytes = 1000/g' ~/.banksy/config/config.toml
cat $HOME/.banksy/config/genesis.json | jq '.app_state["gov"]["params"]["voting_period"]="45s"' > $HOME/.banksy/config/tmp_genesis.json && mv $HOME/.banksy/config/tmp_genesis.json $HOME/.banksy/config/genesis.json

# Start the node (remove the --pruning=nothing flag if historical queries are not needed)
centaurid start --pruning=nothing  --minimum-gas-prices=0stake