# composable-testnet

Cosmos testnet with IBC-v7 and wasm client enable.

run `testnode.sh` first then run `upload_contracts.sh` to upload wasm client.

NOTE: at the end of `testnode.sh`, there is a sed command that change max request body so that we can upload the light client. default is for `macos` change if you use `linux`
