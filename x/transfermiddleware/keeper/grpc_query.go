package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
)

func (k Keeper) ParaTokenInfo(c context.Context, req *types.QueryParaTokenInfoRequest) (*types.QueryParaTokenInfoResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	info := k.GetParachainIBCTokenInfo(ctx, req.NativeDenom)

	return &types.QueryParaTokenInfoResponse{
		IbcDenom:    info.IbcDenom,
		NativeDenom: info.NativeDenom,
		ChannelId:   info.ChannelId,
	}, nil
}

func (k Keeper) EscrowAddress(_ context.Context, req *types.QueryEscrowAddressRequest) (*types.QueryEscrowAddressResponse, error) {
	escrowAddress := transfertypes.GetEscrowAddress(transfertypes.PortID, req.ChannelId)

	return &types.QueryEscrowAddressResponse{
		EscrowAddress: escrowAddress.String(),
	}, nil
}
