#!/bin/bash
# Run this script to quickly install, setup, and run the current version of juno without docker.
# ./scripts/test_node.sh [clean|c]

KEY="test"
CHAINID="centaurid-t1"
MONIKER="localcentaurid"
KEYALGO="secp256k1"
KEYRING="test"
LOGL="info"

picad config keyring-backend $KEYRING
picad config chain-id $CHAINID

command -v picad > /dev/null 2>&1 || { echo >&2 "centaurid command not found. Ensure this is setup / properly installed in your GOPATH."; exit 1; }
command -v jq > /dev/null 2>&1 || { echo >&2 "jq not installed. More info: https://stedolan.github.io/jq/download/"; exit 1; }

from_scratch () {

  make install

  # remove existing daemon.
  rm -rf ~/.banksy/*

  # juno1efd63aw40lxf3n4mhf7dzhjkr453axurv2zdzk
  echo "decorate bright ozone fork gallery riot bus exhaust worth way bone indoor calm squirrel merry zero scheme cotton until shop any excess stage laundry" | picad keys add $KEY --keyring-backend $KEYRING --algo $KEYALGO --recover
  # juno1hj5fveer5cjtn4wd6wstzugjfdxzl0xps73ftl
  echo "wealth flavor believe regret funny network recall kiss grape useless pepper cram hint member few certain unveil rather brick bargain curious require crowd raise" | picad keys add myaccount --keyring-backend $KEYRING --algo $KEYALGO --recover

  picad init $MONIKER --chain-id $CHAINID

  # Function updates the config based on a jq argument as a string
  update_test_genesis () {
    # update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
    cat $HOME/.banksy/config/genesis.json | jq "$1" > $HOME/.banksy/config/tmp_genesis.json && mv $HOME/.banksy/config/tmp_genesis.json $HOME/.banksy/config/genesis.json
  }

  # Set gas limit in genesis
  update_test_genesis '.consensus_params["block"]["max_gas"]="100000000"'
  update_test_genesis '.app_state["gov"]["params"]["voting_period"]="45s"'

  update_test_genesis '.app_state["staking"]["params"]["bond_denom"]="stake"'
  #update_test_genesis '.app_state["bank"]["params"]["send_enabled"]=[{"denom": "stake","enabled": true}]'
  # update_test_genesis '.app_state["staking"]["params"]["min_commission_rate"]="0.100000000000000000"' # sdk 46 only

  update_test_genesis '.app_state["mint"]["params"]["mint_denom"]="stake"'
  update_test_genesis '.app_state["gov"]["deposit_params"]["min_deposit"]=[{"denom": "stake","amount": "1000000"}]'
  update_test_genesis '.app_state["crisis"]["constant_fee"]={"denom": "stake","amount": "1000"}'

  update_test_genesis '.app_state["tokenfactory"]["params"]["denom_creation_fee"]=[{"denom":"stake","amount":"100"}]'

  update_test_genesis '.app_state["feeshare"]["params"]["allowed_denoms"]=["stake"]'

  # Allocate genesis accounts
  picad add-genesis-account $KEY 10000000000000000000stake,100000000000000utest --keyring-backend $KEYRING
  picad add-genesis-account myaccount 1000000000stake,100000000000000utest --keyring-backend $KEYRING

  picad gentx $KEY 10000000000000000000stake --keyring-backend $KEYRING --chain-id $CHAINID

  # Collect genesis tx
  picad collect-gentxs

  # Run this to ensure junorything worked and that the genesis file is setup correctly
  picad validate-genesis
}


if [ $# -eq 1 ] && [ $1 == "clean" ] || [ $1 == "c" ]; then
  echo "Starting from a clean state"
  from_scratch
fi

echo "Starting node..."

# Opens the RPC endpoint to outside connections
sed -i '/laddr = "tcp:\/\/127.0.0.1:26657"/c\laddr = "tcp:\/\/0.0.0.0:26657"' ~/.banksy/config/config.toml
sed -i 's/cors_allowed_origins = \[\]/cors_allowed_origins = \["\*"\]/g' ~/.banksy/config/config.toml
sed -i 's/enable = false/enable = true/g' ~/.banksy/config/app.toml
sed -i '/address = "tcp:\/\/0.0.0.0:1317"/c\address = "tcp:\/\/0.0.0.0:1318"' ~/.banksy/config/app.toml

picad config node tcp://0.0.0.0:2241
picad start --pruning=nothing  --minimum-gas-prices=0stake --p2p.laddr tcp://0.0.0.0:2240 --rpc.laddr tcp://0.0.0.0:2241 --grpc.address 0.0.0.0:2242 --grpc-web.address 0.0.0.0:2243

#MEMO='{"forward":{"receiver":"cosmos18p5cs3z0q68hq7q0d8tr8kp3ldnqkx2fx3f88w","port":"transfer","channel":"channel-0","timeout":600000000000,"retries":0,"next":"{}"}'
#hermes --config scripts/relayer_hermes/config_compose_gaia.toml create channel --a-chain picad-t1 --b-chain gaiad-t1 --a-port transfer --b-port transfer --new-client-connection --yes
#picad tx ibc-transfer transfer transfer channel-0 cosmos1alc8mjana7ssgeyffvlfza08gu6rtav8rmj6nv 10000000stake --from myaccount --keyring-backend test --chain-id picad-t1 --yes --fees 5000stake