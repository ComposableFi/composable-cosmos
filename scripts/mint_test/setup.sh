#!/bin/bash
# Run this script to quickly install, setup, and run the current version of juno without docker.
# Test address cosmos1hce8cea32gjg9eaqzxj02jrgc6m6q59wly4zpm

CHANNEL_ID="channel-0"

ESCROW_ADDRESS=$(picad q transfermiddleware escrow-address channel-0)

hermes --config scripts/relayer_hermes/config_compose_gaia.toml create channel --a-chain centaurid-t1 --b-chain gaiad-t1 --a-port transfer --b-port transfer --new-client-connection --yes

gaiad tx ibc-transfer transfer transfer channel-0 "$ESCROW_ADDRESS" 1000000000stake --from gnad --keyring-backend test --chain-id gaiad-t1 --yes --fees 5000stake
sleep 20
balancesEscrowAdress = $(picad query bank balances $ESCROW_ADDRESS)

picad

picad tx ibc-transfer transfer transfer channel-0 cosmos1hce8cea32gjg9eaqzxj02jrgc6m6q59wly4zpm

