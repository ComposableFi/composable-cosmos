package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/x/mint/types"
)

// RegisterInvariants registers the mint module invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "positive-annual-provisions",
		AnnualProvisionsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "inflation-bounds",
		InflationBoundsInvariant(k))
	ir.RegisterRoute(types.ModuleName, "params-validity",
		ParamsValidityInvariant(k))
	ir.RegisterRoute(types.ModuleName, "token-limits",
		TokenLimitsInvariant(k))
}

// AnnualProvisionsInvariant checks that annual provisions are always positive
func AnnualProvisionsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		minter := k.GetMinter(ctx)
		if minter.AnnualProvisions.IsNegative() {
			return sdk.FormatInvariant(
				types.ModuleName, "annual-provisions",
				"annual provisions cannot be negative",
			), true
		}
		return "", false
	}
}

// InflationBoundsInvariant checks that inflation rate stays within bounds
func InflationBoundsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		minter := k.GetMinter(ctx)

		if minter.Inflation.IsNegative() {
			return sdk.FormatInvariant(
				types.ModuleName, "inflation-bounds",
				"inflation rate cannot be negative",
			), true
		}

		if minter.Inflation.GT(sdk.NewDec(1)) {
			return sdk.FormatInvariant(
				types.ModuleName, "inflation-bounds",
				"inflation rate cannot exceed 100%",
			), true
		}

		return "", false
	}
}

// ParamsValidityInvariant checks that mint parameters are valid
func ParamsValidityInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		params := k.GetParams(ctx)

		if params.BlocksPerYear == 0 {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"blocks per year must be positive",
			), true
		}

		if params.MintDenom == "" {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"mint denom cannot be empty",
			), true
		}

		if params.GoalBonded.IsNegative() {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"goal bonded ratio cannot be negative",
			), true
		}

		if params.GoalBonded.GT(sdk.NewDec(1)) {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"goal bonded ratio cannot exceed 1",
			), true
		}

		if params.InflationRateChange.IsNegative() {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"inflation rate change cannot be negative",
			), true
		}

		if params.MaxTokenPerYear.IsNegative() {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"max token per year cannot be negative",
			), true
		}

		if params.MinTokenPerYear.IsNegative() {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"min token per year cannot be negative",
			), true
		}

		if params.MinTokenPerYear.GT(params.MaxTokenPerYear) {
			return sdk.FormatInvariant(
				types.ModuleName, "params-validity",
				"min token per year cannot exceed max token per year",
			), true
		}

		return "", false
	}
}

// TokenLimitsInvariant checks that annual provisions stay within min/max token limits
func TokenLimitsInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		minter := k.GetMinter(ctx)
		params := k.GetParams(ctx)

		// Convert Int to Dec for comparison
		minTokenDec := sdk.NewDecFromInt(params.MinTokenPerYear)
		maxTokenDec := sdk.NewDecFromInt(params.MaxTokenPerYear)

		if minTokenDec.IsPositive() && minter.AnnualProvisions.LT(minTokenDec) {
			return sdk.FormatInvariant(
				types.ModuleName, "token-limits",
				fmt.Sprintf("annual provisions %s is below minimum %s", minter.AnnualProvisions, minTokenDec),
			), true
		}

		if maxTokenDec.IsPositive() && minter.AnnualProvisions.GT(maxTokenDec) {
			return sdk.FormatInvariant(
				types.ModuleName, "token-limits",
				fmt.Sprintf("annual provisions %s is above maximum %s", minter.AnnualProvisions, maxTokenDec),
			), true
		}

		return "", false
	}
}
