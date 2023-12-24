package v5_1_0

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/composable/v6/app/keepers"
	rateLimitKeeper "github.com/notional-labs/composable/v6/x/ratelimit/keeper"
	"github.com/notional-labs/composable/v6/x/ratelimit/types"
)

const uosmo = "ibc/47BD209179859CDE4A2806763D7189B6E6FE13A17880FE2B42DE1E6C1E329E23"

func RunForkLogic(ctx sdk.Context, keepers *keepers.AppKeepers) {
	ctx.Logger().Info("Applying v5_1_0 upgrade" +
		"Fix Rate Limit With Osmosis Token",
	)

	FixRateLimit(ctx, &keepers.RatelimitKeeper)
}

func FixRateLimit(ctx sdk.Context, rlKeeper *rateLimitKeeper.Keeper) {
	uosmoRateLimit, found := rlKeeper.GetRateLimit(ctx, uosmo, "channel-2")
	if !found {
		channelValue := rlKeeper.GetChannelValue(ctx, uosmo)
		// Create and store the rate limit object
		path := types.Path{
			Denom:     uosmo,
			ChannelID: "channel-2",
		}
		quota := types.Quota{
			MaxPercentSend: sdkmath.NewInt(30),
			MaxPercentRecv: sdkmath.NewInt(30),
			DurationHours:  24,
		}
		flow := types.Flow{
			Inflow:       math.ZeroInt(),
			Outflow:      math.ZeroInt(),
			ChannelValue: channelValue,
		}
		uosmoRateLimit = types.RateLimit{
			Path:               &path,
			Quota:              &quota,
			Flow:               &flow,
			MinRateLimitAmount: sdkmath.NewInt(1), // decimal 6
		}
		rlKeeper.SetRateLimit(ctx, uosmoRateLimit)
	} else {
		uosmoRateLimit.MinRateLimitAmount = sdkmath.NewInt(1)
		rlKeeper.SetRateLimit(ctx, uosmoRateLimit)
	}

	// double check
	allRateLiit := rlKeeper.GetAllRateLimits(ctx)
	for _, ratelimit := range allRateLiit {
		if ratelimit.MinRateLimitAmount.IsNil() {
			ratelimit.MinRateLimitAmount = sdkmath.NewInt(1)
			rlKeeper.SetRateLimit(ctx, ratelimit)
		}
	}
}
