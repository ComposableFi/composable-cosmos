package prepare

import (
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	txboundarykeeper "github.com/notional-labs/composable/v6/x/tx-boundary/keeper"
)

func PrepareProposalHandler(
	txConfig client.TxConfig,
	cdc codec.BinaryCodec,
	txboundaryKeeper txboundarykeeper.Keeper,
) sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req abci.RequestPrepareProposal) abci.ResponsePrepareProposal {
		filterChangeValSetMsgs := FilterChangeValSetMsgs(ctx, cdc, txConfig.TxDecoder(), req.Txs)
		if req.Height%5 == 0 {
			return abci.ResponsePrepareProposal{Txs: req.Txs}
		}

		return abci.ResponsePrepareProposal{Txs: filterChangeValSetMsgs}
	}
}
