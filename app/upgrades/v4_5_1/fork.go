package v4_5_1

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/centauri/v4/app/keepers"
	rateLimitKeeper "github.com/notional-labs/centauri/v4/x/ratelimit/keeper"
)

func RunForkLogic(ctx sdk.Context, keepers *keepers.AppKeepers) {
	ctx.Logger().Info("Applying v5 upgrade" +
		"Remove Rate Limit",
	)

	RemoveRateLimit(ctx, &keepers.RatelimitKeeper)
}

func RemoveRateLimit(ctx sdk.Context, rlKeeper *rateLimitKeeper.Keeper) {
	// Get all current rate limit
	rateLimits := rlKeeper.GetAllRateLimits(ctx)
	// Remove Rate limit
	for _, rateLimit := range rateLimits {
		err := rlKeeper.RemoveRateLimit(ctx, rateLimit.Path.Denom, rateLimit.Path.ChannelId)
		if err != nil {
			panic(err)
		}
	}
}
