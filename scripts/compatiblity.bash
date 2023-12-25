# compatibility check
# we should currently be comparing only to v6.3.1
# git checkout v6.3.1 

go install ./...
centaurid init
wget https://snapshots.polkachu.com/snapshots/composable/composable_2959777.tar.lz4
wget https://raw.githubusercontent.com/notional-labs/composable-networks/main/mainnet/genesis.json
mv genesis.json  ~/.banksy/config/genesis.json
mv composable_2959777.tar.lz4 ~/.banksy/composable_2959777.tar.lz4
cd ~/.banksy
lz4 -d composable_2959777.tar.lz4 | tar -xvf -
timeout 5m centaurid start --p2p.seeds "ebc272824924ea1a27ea3183dd0b9ba713494f83@composable-mainnet-seed.autostake.com:26976,20e1000e88125698264454a884812746c2eb4807@seeds.lavenderfive.com:22256,d2362ebcdd562500ac8c4cfa2214a89ad811033c@seeds.whispernode.com:22256"
