#!/bin/bash

# the upgrade is a fork, "true" otherwise
FORK=${FORK:-"false"}

OLD_VERSION=v6.4.3
UPGRADE_WAIT=${UPGRADE_WAIT:-20}
HOME=mytestnet
ROOT=$(pwd)
DENOM=stake
CHAIN_ID1=test-1
SOFTWARE_UPGRADE_NAME="v6_4_8"
ADDITIONAL_PRE_SCRIPTS="./scripts/upgrade/v_6_4_8/pre-script-ibc-transfer.sh"
ADDITIONAL_AFTER_SCRIPTS="./scripts/upgrade/v_6_4_8/post-script-ibc-transfer.sh"
SETUP_RELAYER_SCRIPTS="./scripts/relayer/relayer-init.sh"

CHAIN_DIR=$(pwd)/data
CHAINID_1=test-1

SLEEP_TIME=1

if [[ "$FORK" == "true" ]]; then
    export PICA_HALT_HEIGHT=20
fi

# underscore so that go tool will not take gocache into account
mkdir -p _build/gocache
export GOMODCACHE=$ROOT/_build/gocache

# install old binary if not exist
if [ ! -f "_build/$OLD_VERSION.zip" ] &> /dev/null
then
    mkdir -p _build/old
    wget -c "https://github.com/ComposableFi/composable-cosmos/archive/refs/tags/${OLD_VERSION}.zip" -O _build/${OLD_VERSION}.zip
    unzip _build/${OLD_VERSION}.zip -d _build
fi

# reinstall old binary
if [ $# -eq 1 ] && [ $1 == "--reinstall-old" ] || ! command -v _build/old/centaurid &> /dev/null; then
    cd ./_build/composable-cosmos-${OLD_VERSION:1}
    GOBIN="$ROOT/_build/old" go install -mod=readonly ./...
    cd ../..
fi

# install new binary
if ! command -v _build/new/picad &> /dev/null
then
    mkdir -p _build/new
    GOBIN="$ROOT/_build/new" go install -mod=readonly ./...
fi

# run old node
echo "running old node"
bash scripts/two-testnodes.sh _build/old/centaurid

sleep 5 # wait for 2 node to start 

#setup ibc between 2 nodes
echo "setting up ibc"
source $SETUP_RELAYER_SCRIPTS

# Transfer from chain 1 to chain 2 and return balance of sender on chain 1
echo "executing additional pre scripts from $ADDITIONAL_PRE_SCRIPTS"
bash ./scripts/upgrade/v_6_4_8/pre-script-ibc-transfer.sh

run_fork () {
    echo "forking"

    while true; do
        BLOCK_HEIGHT=$(./_build/old/centaurid status | jq '.SyncInfo.latest_block_height' -r)
        # if BLOCK_HEIGHT is not empty
        if [ ! -z "$BLOCK_HEIGHT" ]; then
            echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
            sleep 10
        else
            echo "BLOCK_HEIGHT is empty, forking"
            break
        fi
    done
}

run_upgrade () {
    echo "start upgrading chain-1"

    # Get upgrade height, 12 block after (6s)
    STATUS_INFO=($(./_build/old/centaurid status --home $CHAIN_DIR/$CHAINID_1 | jq -r '.NodeInfo.network,.SyncInfo.latest_block_height'))
    UPGRADE_HEIGHT=$((STATUS_INFO[1] + 20))
    echo "UPGRADE_HEIGHT = $UPGRADE_HEIGHT"

    tar -cf ./_build/new/picad.tar -C ./_build/new picad
    SUM=$(shasum -a 256 ./_build/new/picad.tar | cut -d ' ' -f1)
    UPGRADE_INFO=$(jq -n '
    {
        "binaries": {
            "linux/amd64": "file://'$(pwd)'/_build/new/picad.tar?checksum=sha256:'"$SUM"'",
        }
    }')

    ./_build/old/centaurid keys list --home $CHAIN_DIR/$CHAINID_1 --keyring-backend test

    ./_build/old/centaurid tx gov submit-legacy-proposal software-upgrade "$SOFTWARE_UPGRADE_NAME" --upgrade-height $UPGRADE_HEIGHT --upgrade-info "$UPGRADE_INFO" --title "upgrade" --description "upgrade"  --from val1 --keyring-backend test --chain-id $CHAIN_ID1 --home $CHAIN_DIR/$CHAINID_1 --node tcp://localhost:16657 --output=json -y > /dev/null

    sleep $SLEEP_TIME

    ./_build/old/centaurid tx gov deposit 1 "20000000${DENOM}" --from val1 --keyring-backend test --chain-id $CHAIN_ID1 --home $CHAIN_DIR/$CHAINID_1 --node tcp://localhost:16657 --output=json -y > /dev/null

    sleep $SLEEP_TIME

    ./_build/old/centaurid tx gov vote 1 yes --from val1 --keyring-backend test --chain-id $CHAIN_ID1 --home $CHAIN_DIR/$CHAINID_1 --node tcp://localhost:16657 --output=json -y > /dev/null

    sleep $SLEEP_TIME

    # determine block_height to halt
    while true; do
        BLOCK_HEIGHT=$(./_build/old/centaurid status --home $CHAIN_DIR/$CHAINID_1 | jq '.SyncInfo.latest_block_height' -r)
        if [ $BLOCK_HEIGHT = "$UPGRADE_HEIGHT" ]; then
            # only kill the first centaurid
            echo "BLOCK HEIGHT = $UPGRADE_HEIGHT REACHED, KILLING OLD ONE"
            pkill -o centaurid
            break
        else
            ./_build/old/centaurid q gov proposal 1 --home $CHAIN_DIR/$CHAINID_1 --output=json | jq ".status"
            echo "BLOCK_HEIGHT = $BLOCK_HEIGHT"
            sleep 1 
        fi
    done
}

# if FORK = true
if [[ "$FORK" == "true" ]]; then
    run_fork
    unset PICA_HALT_HEIGHT
else
    run_upgrade
fi

sleep 1

echo ""
echo "#########################################################"
echo "# UPGRADE CHAIN 1 SUCCESSFUL #"
echo "#########################################################"
echo ""

# run new node
CONTINUE="true" bash scripts/two-testnodes.sh _build/new/picad

sleep 5

echo "executing additional post scripts from $ADDITIONAL_AFTER_SCRIPTS"
bash ./scripts/upgrade/v_6_4_8/post-script-ibc-transfer.sh

echo ""
echo "#########################################################"
echo "# INTERCHAIN TEST SUCCESSFUL #"
echo "#########################################################"
echo ""