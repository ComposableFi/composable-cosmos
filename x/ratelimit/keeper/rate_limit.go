package keeper

import (
	"strings"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/notional-labs/centauri/v3/x/ratelimit/types"
)

// The total value on a given path (aka, the denominator in the percentage calculation)
// is the total supply of the given denom
func (k Keeper) GetChannelValue(ctx sdk.Context, denom string) math.Int {
	return k.bankKeeper.GetSupply(ctx, denom).Amount
}

// If the rate limit is exceeded or the denom is blacklisted, we emit an event
func EmitTransferDeniedEvent(ctx sdk.Context, reason, denom, channelId string, direction types.PacketDirection, amount math.Int, err error) {
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTransferDenied,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyReason, reason),
			sdk.NewAttribute(types.AttributeKeyAction, strings.ToLower(direction.String())), // packet_send or packet_recv
			sdk.NewAttribute(types.AttributeKeyDenom, denom),
			sdk.NewAttribute(types.AttributeKeyChannel, channelId),
			sdk.NewAttribute(types.AttributeKeyAmount, amount.String()),
			sdk.NewAttribute(types.AttributeKeyError, err.Error()),
		),
	)
}

// Adds an amount to the flow in either the SEND or RECV direction
func (k Keeper) UpdateFlow(rateLimit types.RateLimit, direction types.PacketDirection, amount math.Int) error {
	switch direction {
	case types.PACKET_SEND:
		return rateLimit.Flow.AddOutflow(amount, *rateLimit.Quota)
	case types.PACKET_RECV:
		return rateLimit.Flow.AddInflow(amount, *rateLimit.Quota)
	default:
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid packet direction (%s)", direction.String())
	}
}
