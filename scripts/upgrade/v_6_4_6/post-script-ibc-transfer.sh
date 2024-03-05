#!/bin/bash

echo ""
echo "#################################################"
echo "# post-script: IBC Transfer #"
echo "#################################################"
echo ""

NEW_BINARY=picad
BINARY=centaurid
CHAIN_DIR=$(pwd)/data

AMOUNT_TO_DELEGATE=10000000000
UPICA_DENOM=stake
WALLET_1=$($NEW_BINARY keys show wallet1 -a --keyring-backend test --home $CHAIN_DIR/test-1)
WALLET_2=$($BINARY keys show wallet2 -a --keyring-backend test --home $CHAIN_DIR/test-2)

ACCOUNT_BALANCE_OF_WALLET1=$($NEW_BINARY q bank balances $WALLET_1 --chain-id test-1 --node tcp://localhost:16657 -o json | jq -r '.balances[0].amount')
echo "CHECK BALANCE POST UPGRADE: $ACCOUNT_BALANCE_OF_WALLET1"

echo "Sending tokens from validator wallet on test-1 to validator wallet on test-2"
IBC_TRANSFER=$($NEW_BINARY tx ibc-transfer transfer transfer channel-0 $WALLET_2 $AMOUNT_TO_DELEGATE$UPICA_DENOM --chain-id test-1 --from $WALLET_1 --home $CHAIN_DIR/test-1 --fees 60000$UPICA_DENOM --node tcp://localhost:16657 --keyring-backend test  -y -o json | jq -r '.raw_log' )

if [[ "$IBC_TRANSFER" == "failed to execute message"* ]]; then
    echo "Error: IBC transfer failed, with error: $IBC_TRANSFER"
    exit 1
fi

sleep 5

ACCOUNT_BALANCE=$($NEW_BINARY q bank balances $WALLET_1 --chain-id test-1 --node tcp://localhost:16657 -o json | jq -r '.balances[0].amount')
echo "CHECK BALANCE POST UPGRADE AFTER SEND SECOND TIME: $ACCOUNT_BALANCE"

exit 1

echo ""
echo "#########################################################"
echo "# Success: post-script: IBC Transfer #"
echo "#########################################################"
echo ""