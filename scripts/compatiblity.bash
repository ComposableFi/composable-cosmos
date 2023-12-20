# compatibility check
# we should currently be comparing only to v6.3.1
# git checkout v6.3.1 

go install ./...
centaurid init
aria2c https://snapshots.polkachu.com/snapshots/composable/composable_2959777.tar.lz4
aria2c https://raw.githubusercontent.com/notional-labs/composable-networks/main/mainnet/genesis.json
mv genesis.json  ~/.banksy/config/genesis.json
mv composable_2959777.tar.lz4 ~/.banksy/composable_2959777.tar.lz4
cd ~/.banksy
lz4 -d composable_2959777.tar.lz4 | tar -xvf -
centaurid start
