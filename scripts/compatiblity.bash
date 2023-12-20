aria2c https://snapshots.polkachu.com/snapshots/composable/composable_2959777.tar.lz4
mv composable_2959777.tar.lz4 ~/.banksy/composable_2959777.tar.lz4
cd ~/.banksy
lz4 -d composable_2959777.tar.lz4 | tar -xvf -
rm composable_2959777.tar.lz4
centaurid start
