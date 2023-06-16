#!/bin/sh

CHAIN_ID=localcentauri
CENTAURI_HOME=$HOME/.banksy
CONFIG_FOLDER=$CENTAURI_HOME/config
MONIKER=val
STATE='false'

# val - centauri1jxa3ksucx7ter57xyuczvmk6qkeqmqvj37g237
MNEMONIC="blame tube add leopard fire next exercise evoke young team payment senior know estate mandate negative actual aware slab drive celery elevator burden utility"

while getopts s flag
do
    case "${flag}" in
        s) STATE='true';;
    esac
done

install_prerequisites () {
    apk add dasel
}

edit_genesis () {

    GENESIS=$CONFIG_FOLDER/genesis.json

    # Update staking module
    dasel put string -f $GENESIS '.app_state.staking.params.bond_denom' 'ppica'
    dasel put string -f $GENESIS '.app_state.staking.params.unbonding_time' '240s'

    # Update crisis module
    dasel put string -f $GENESIS '.app_state.crisis.constant_fee.denom' 'ppica'

    # Udpate gov module
    dasel put string -f $GENESIS '.app_state.gov.voting_params.voting_period' '60s'
    dasel put string -f $GENESIS '.app_state.gov.deposit_params.min_deposit.[0].denom' 'ppica'

    # Update wasm permission (Nobody or Everybody)
    dasel put string -f $GENESIS '.app_state.wasm.params.code_upload_access.permission' "Everybody"
}

add_genesis_accounts () {

    centaurid add-genesis-account centauri1jxa3ksucx7ter57xyuczvmk6qkeqmqvjrxxdj3 1000000000000000ppica --home $CENTAURI_HOME # val
    centaurid add-genesis-account centauri1cyyzpxplxdzkeea7kwsydadg87357qnamvg3y3 1000000000000000ppica --home $CENTAURI_HOME # lo-test1
    centaurid add-genesis-account centauri18s5lynnmx37hq4wlrw9gdn68sg2uxp5ry85k7d 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri1qwexv7c6sm95lwhzn9027vyu2ccneaqapystyu 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri14hcxlnwlqtq75ttaxf674vk6mafspg8xzedlxs 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri12rr534cer5c0vj53eq4y32lcwguyy7nnp63sc2 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri1nt33cjd5auzh36syym6azgc8tve0jlvknz7jqp 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri10qfrpash5g2vk3hppvu45x0g860czur89c2g5w 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri1f4tvsdukfwh6s9swrc24gkuz23tp8pd345acmu 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri1myv43sqgnj5sm4zl98ftl45af9cfzk7nmrc7jk 1000000000000000ppica --home $CENTAURI_HOME
    centaurid add-genesis-account centauri14gs9zqh8m49yy9kscjqu9h72exyf295a9ey66h 1000000000000000ppica --home $CENTAURI_HOME # lo-test10

    echo $MNEMONIC | centaurid keys add $MONIKER --recover --keyring-backend=test --home $CENTAURI_HOME
    centaurid gentx $MONIKER 50000000000000ppica --keyring-backend=test --chain-id=$CHAIN_ID --home $CENTAURI_HOME

    centaurid collect-gentxs --home $CENTAURI_HOME
}

edit_config () {
    # Remove seeds
    dasel put string -f $CONFIG_FOLDER/config.toml '.p2p.seeds' ''

    # Expose the rpc
    dasel put string -f $CONFIG_FOLDER/config.toml '.rpc.laddr' "tcp://0.0.0.0:26657"
}



if [[ ! -d $CONFIG_FOLDER ]]
then
    echo $MNEMONIC | centaurid init -o --chain-id=$CHAIN_ID --home $CENTAURI_HOME --recover $MONIKER
    install_prerequisites
    edit_genesis
    add_genesis_accounts
    edit_config
fi

centaurid start --home $CENTAURI_HOME &
# killall centaurid

wait