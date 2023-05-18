# Composable Testnet 2 Genesis Transaction

## Binary
Compile from source (go 1.20 recommended):
```
git clone https://github.com/notional-labs/composable-testnet
cd composable-testnet 
git checkout v2.3.0 
make install
```

## Init
```
banksyd init NODE_NAME --chain-id banksy-testnet-2
wget -O ~/.banksy/config/genesis.json https://raw.githubusercontent.com/notional-labs/composable-testnet/main/networks/testnet-2/pregenesis.json
banksyd config chain-id banksy-testnet-2
```

## Keys
Generate a new key:
```
banksyd keys add KEYNAME 
```
Recover existing key:
```
banksyd keys add KEYNAME --recover
```

## Creating genesis transaction
Step-by-step guide:
```
banksyd add-genesis-account KEYNAME 1000000upica
banksyd gentx KEYNAME 1000000upica \
--moniker="" \
--identity="" \
--details="" \
--website="" \
--security-contact="" \
--chain-id banksy-testnet-2
```
The output will look like this: 
```
Genesis transaction written to "~/.banksy/config/gentx/gentx-799d25f37dc6c68a549abbcd98e73127ac60d492.json"
```
Fork the repo and create a pull request with your gentx-XXX.json moved to this directory: https://github.com/notional-labs/composable-testnet/tree/main/networks/testnet-2/gentxs
Remember to change the file name to your validator name `gentx-YOURNAME.json`
Example:
```
mv ~/.banksy/config/gentx/gentx-XXX.json ~/composable-testnet/networks/testnet-2/gentxs/gentx-YOURNAME.json
```
