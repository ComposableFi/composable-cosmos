#!/bin/bash

# this bash will prepare cosmosvisor to the build folder so that it can perform upgrade
# this script is supposed to be run by Makefile

# These fields should be fetched automatically in the future
# Need to do more upgrade to see upgrade patterns
OLD_VERSION=v6.4.x
# this command will retrieve the folder with the largest number in format v<number>
SOFTWARE_UPGRADE_NAME="v7"
BUILDDIR=$1
TESTNET_NVAL=$2
TESTNET_CHAINID=$3

# check if BUILDDIR is set
if [ -z "$BUILDDIR" ]; then
    echo "BUILDDIR is not set"
    exit 1
fi

# install old binary if not exist
if [ ! -f "_build/$OLD_VERSION.zip" ] &> /dev/null
then
    mkdir -p _build/old
    # This archive have cmd testnet for upgrade testing
    wget -c "https://github.com/tungle-notional/composable-old/archive/refs/tags/v6.4.x_old.zip" -O _build/${OLD_VERSION}.zip
    unzip _build/${OLD_VERSION}.zip -d _build
fi


if [ ! -f "$BUILDDIR/old/centaurid" ] &> /dev/null
then
    if [ ! "$(docker images -q centauri/centaurid.binary.old 2> /dev/null)" ]; then
        docker build --platform linux/amd64 --no-cache --build-arg source=./_build/composable-cosmos-${OLD_VERSION:1}/ --tag centauri/centaurid.binary.old ./_build/composable-cosmos-${OLD_VERSION:1}
    fi
    docker create --platform linux/amd64 --name old-temp centauri/centaurid.binary.old:latest
    mkdir -p $BUILDDIR/old
    docker cp old-temp:/bin/centaurid $BUILDDIR/old/
    docker rm old-temp
fi

echo "init-files"
# prepare cosmovisor config in TESTNET_NVAL nodes
if [ ! -f "$BUILDDIR/node0/centaurid/config/genesis.json" ]; then docker run --rm \
    -v $BUILDDIR:/centaurid:Z \
    --platform linux/amd64 \
    --entrypoint /centaurid/old/centaurid \
    centauri/centaurid-upgrade-env testnet init-files --v $TESTNET_NVAL --chain-id $TESTNET_CHAINID -o . --starting-ip-address 192.168.0.2 --minimum-gas-prices "0stake" --node-daemon-home centaurid --keyring-backend=test --home=temp; \
fi

for (( i=0; i<$TESTNET_NVAL; i++ )); do
    CURRENT=$BUILDDIR/node$i/centaurid
    echo "Change voting_period"
    # change gov params voting_period
    jq '.app_state["gov"]["params"]["voting_period"] = "50s"' $CURRENT/config/genesis.json > $CURRENT/config/genesis.json.tmp && mv $CURRENT/config/genesis.json.tmp $CURRENT/config/genesis.json

    docker run --rm \
        -v $BUILDDIR:/centaurid:Z \
        -e DAEMON_HOME=/centaurid/node$i/centaurid \
        -e DAEMON_NAME=centaurid \
        -e DAEMON_RESTART_AFTER_UPGRADE=true \
        --entrypoint /centaurid/cosmovisor \
        --platform linux/amd64 \
        centauri/centaurid-upgrade-env init /centaurid/old/centaurid
    mkdir -p $CURRENT/cosmovisor/upgrades/$SOFTWARE_UPGRADE_NAME/bin
    cp $BUILDDIR/centaurid $CURRENT/cosmovisor/upgrades/$SOFTWARE_UPGRADE_NAME/bin
    touch $CURRENT/cosmovisor/upgrades/$SOFTWARE_UPGRADE_NAME/upgrade-info.json
done