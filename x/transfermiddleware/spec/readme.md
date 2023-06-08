# Transfer middleware
This module keep Centauri and Picasso have shared total supply
## Definitions

Middleware: A self-contained module that sits between core IBC and an underlying IBC application during packet execution. All messages between core IBC and underlying application must flow through middleware, which may perform its own custom logic. 

In this case, the IBC application is `transfer` module in ibc-go.

## How is it working 
As we want a shared total supply between Centauri and Picasso, token transfers from Picasso to Centauri via IBC will not have the `ibc/` denom prefix. Instead, the native token in Picasso, PICA, should be used.

We have introduced an IBC middleware module to lock `ibc/` tokens and mint native tokens. This module handles middleware processes when executing the ICS26 implementation of IBC transfer, including SendPacket, OnRecvPacket, OnAcknowledgementPacket, and OnTimeoutPacket.

#### Transfer and SendPacket
![](https://hackmd.io/_uploads/Hy3dFx4M2.png)

Scope:
 - Coverage source zone, sink zone logic in IBC transfer. (about this : https://github.com/cosmos/ibc-go/blob/main/modules/apps/transfer/keeper/relay.go#L22)
 - Ensure all native and ibced token will be burned when receive successful ack packet.

Current process :  https://github.com/notional-labs/composable-testnet/pull/25

#### OnRecvPacket
![](https://hackmd.io/_uploads/BJAL7BNfn.png)

Scope:
 - Coverage source zone, sink zone logic in IBC transfer. (about this : https://github.com/cosmos/ibc-go/blob/main/modules/apps/transfer/keeper/relay.go#L22)

#### OnAcknowledgementPacket


#### OnTimeoutPacket


## Burn and mint PICA token when launch chain and when IBC connected :
Before the bridge launch, there will be approximately 1 billion PICA tokens on the Banksy chain and 9 billion on the Picasso chain. The supply of PICA is fixed, and the tokens cannot flow between the two chains yet.

After the bridge launch, PICA tokens will be able to flow between the two chains, so the amount on each chain will change, but the total supply will remain the same.

There's a problem with escrow address in IBC module, but we can handle this by a burning and IBC transfer process. We need to test this feature in testnet-2 before launch mainnet.
