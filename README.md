# Composable-Cosmos

Picasso is a blockchain built with the cosmos-sdk that uses [IBC](github.com/cosmos/ibc-go) to securely connect chains like:

* Polkadot
* Solana
* Ethereum

...to the Cosmos ecosystem, and to one another.  We use [CosmWasm](github.com/CosmWasm/wasmd) to provide a contract platform that allows developers to deploy smart contracts written in rust to both Cosmos and Polkadot, and for those contracts to interact with one another.  On Solana and Ethereum, contracts deployed on both chains allow for the full set of IBC security and communication features to be used, so that developers can seamlessly integrate applications between ecosystems.

As we proceed with our journey, front end libraries that make building across the entire set of IBC connected chains will be released, enabling a new class of Web3 applications that knows no boundaries.

To learn more, or just to reach out, check out our [X](https://twitter.com/ComposableFin) or [Discord](https://t.co/eilU0TYTYN).  


## Hardware Recommendation

* Quad core or larger amd64 CPU
* 64GB+ RAM
* 1TB+ NVMe Storage

## Quick start

Requires [Go 1.21](https://go.dev/doc/install) or higher.

```bash
make install
centaurid version
```
Then you can run a node with a single command.

```bash
./scripts/testnode.sh   
```

If you have Docker installed, then you can run a local node with a single command.

```bash
docker compose up -d
```
and remove all node running with

```bash
docker compose down
```
