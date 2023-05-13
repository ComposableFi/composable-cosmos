package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
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
