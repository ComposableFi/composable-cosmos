# Transfer Middleware integrated with Packet-Forward-Middleware
## Context
As the requirment of Composable (we want to mint native token that represent ibc-token that was sent from Picasso), the Transfer-middleware had been developed. The ideal of Transfer-middleware is lock the ibc-token that was sent and mint the same amount of native-token to receiver address.

=> Transfer-middleware will effect to packet-forward-middleware, so we need to handle it.

## Problem when packet forward from Picasso to Cosmos chain

When we want to forward a transfer packet from Picasso to a Cosmos chain via Composable, Composable will lock ibc-token from Picasso, mint native token and then send that native token to Cosmos chain.

### Handle transfer package with PFM memo
When we receive a packet that will be forwarded, we will check if that packet was received from Picasso. If then, we will:
- Send the ibc token to escrow address
- Mint native token to the receiver address on Composable
- Disable denom composition. So that, the PFM module will not replace the transfer denom and we can forward the transfer package with the native denom that was minted.

### Handle Ack
When we receive the package with memo for forwarding, we will extract the memo to store the transfer data package. The data will be used to check the ACK. When we receiver the ACK the responsible with the forwarding package, we will handle it.

Successful case:
- The ibc-token that received on Composable will be locked in the escrow address. TransferMiddleware will mint native token and send to receiver
- PFM will create a new transfer msg that responsible with the MEMO that will send the amount of native token just minted from receiver on Composable and send that msg to Cosmos chain

Error Ack case: 
- PFM module will refund the amount of sent token to the sender (Picasso)
- Transfermiddleware module will recovery the amount of token that lock to the escrow address (burn the amount of token in the escros address)

### Handle TimeOut
- PFM module will refund the amount of sent token to the sender (Picasso)
- Transfermiddleware module will recovery the amount of token that lock to the escrow address (burn the amount of token in the escros address)

## Problem when packet forward from Cosmos chain to Picasso
When we want to forward a transfer packet from Cosmos chain to Picasso, there are 2 scenario:
- Cosmos chain token.
- Token that was received from Picasso.

In the first case, every things will active like normal. But it will have a problem with the second case.

### Handle transfer package
Composable after receive a transfer message with a PFM that forward to Picasso, it will check if the transfer token is the parachain token that store in the state. If so, It will:
- send the amount of locked native token in escrow address to receiver on Composable
- burn the amount of token that just unlocked
- unlock the amount of ibc token in escrow address
- transfer the amount of ibc token just unlocked to Picasso
### Handle Ack
Successful case: Then nothing happen.
Error Ack case: 
- PFM module will refund the amount of sent token to the sender (Cosmos chain)
- Transfermiddleware module will recovery the amount of token that lock to the escrow address (mint the amount of token in the escros address)
### Handle TimeOut
- PFM module will refund the amount of sent token to the sender (Cosmos chain)
- Transfermiddleware module will recovery the amount of token that lock to the escrow address (mint the amount of token in the escros address)

