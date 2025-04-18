package mint

import (
	"fmt"
	"runtime/debug"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/composable/v6/x/mint/keeper"
	"github.com/notional-labs/composable/v6/x/mint/types"
)

// BeginBlocker mints new tokens for the previous block.
func BeginBlocker(ctx sdk.Context, k keeper.Keeper, ic types.InflationCalculationFn) {
	defer func() {
		if r := recover(); r != nil {
			// Log the panic and stack trace
			errMsg := fmt.Sprintf("panic in mint BeginBlocker: %v\n%s", r, string(debug.Stack()))
			ctx.Logger().Error(errMsg)

			// Emit an event for monitoring
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeMint,
					sdk.NewAttribute(types.AttributeKeyError, errMsg),
				),
			)
		}
	}()

	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	// fetch stored minter & params
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// recalculate inflation rate
	totalStakingSupply := k.StakingTokenSupply(ctx)
	bondedRatio := k.BondedRatio(ctx)
	minter.Inflation = ic(ctx, minter, params, bondedRatio, totalStakingSupply)
	annualProvisions, err := minter.NextAnnualProvisions(params, totalStakingSupply)
	if err != nil {
		k.Logger(ctx).Error("failed to calculate annual provisions", "error", err)
		return
	}
	minter.AnnualProvisions = annualProvisions
	k.SetMinter(ctx, minter)

	// calculate how many we would mint, but we dont mint them, we take them from the prefunded account
	mintedCoin, err := minter.BlockProvision(params)
	if err != nil {
		k.Logger(ctx).Error("failed to calculate block provision", "error", err)
		return
	}
	mintedCoins := sdk.NewCoins(mintedCoin)

	// send the minted coins to the fee collector account
	err = k.AddCollectedFees(ctx, mintedCoins)
	if err != nil {
		k.Logger(ctx).Info("Not enough incentive tokens in the mint pool to distribute")
	}

	if mintedCoin.Amount.IsInt64() {
		defer telemetry.ModuleSetGauge(types.ModuleName, float32(mintedCoin.Amount.Int64()), "minted_tokens")
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventTypeMint,
			sdk.NewAttribute(types.AttributeKeyBondedRatio, bondedRatio.String()),
			sdk.NewAttribute(types.AttributeKeyInflation, minter.Inflation.String()),
			sdk.NewAttribute(types.AttributeKeyAnnualProvisions, minter.AnnualProvisions.String()),
			sdk.NewAttribute(sdk.AttributeKeyAmount, mintedCoin.Amount.String()),
		),
	)
}
