package types

import (
	"fmt"
	stdmath "math"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(inflation, annualProvisions sdk.Dec) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdk.Dec) Minter {
	return NewMinter(
		inflation,
		sdk.NewDec(0),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13% per year.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		// Create a new Dec from integer with decimal place at prec
		// CONTRACT: prec <= Precision
		sdk.NewDecWithPrec(InflationRate, Precision),
	)
}

// validate minter
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNil() {
		return fmt.Errorf("minter inflation cannot be nil")
	}
	if minter.AnnualProvisions.IsNil() {
		return fmt.Errorf("minter annual provisions cannot be nil")
	}
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// NextInflationRate returns the new inflation rate for the next hour.
func (m Minter) NextInflationRate(params Params, bondedRatio sdk.Dec, totalStakingSupply math.Int) (sdk.Dec, error) {
	// Validate inputs
	if params.InflationRateChange.IsNil() {
		return sdk.Dec{}, fmt.Errorf("inflation rate change cannot be nil")
	}
	if params.InflationRateChange.IsNegative() {
		return sdk.Dec{}, fmt.Errorf("inflation rate change cannot be negative")
	}
	if bondedRatio.IsNil() {
		return sdk.Dec{}, fmt.Errorf("bonded ratio cannot be nil")
	}
	if bondedRatio.IsNegative() {
		return sdk.Dec{}, fmt.Errorf("bonded ratio cannot be negative")
	}
	if bondedRatio.GT(sdk.OneDec()) {
		return sdk.Dec{}, fmt.Errorf("bonded ratio cannot be greater than 1")
	}
	if totalStakingSupply.IsNil() {
		return sdk.Dec{}, fmt.Errorf("total staking supply cannot be nil")
	}
	if totalStakingSupply.IsNegative() {
		return sdk.Dec{}, fmt.Errorf("staking supply cannot be negative")
	}
	if totalStakingSupply.GT(sdk.NewIntFromUint64(stdmath.MaxUint64)) {
		return sdk.Dec{}, fmt.Errorf("staking supply too large")
	}

	if bondedRatio.IsZero() {
		return sdk.ZeroDec(), nil
	}

	if totalStakingSupply.IsZero() {
		return sdk.ZeroDec(), nil
	}

	// Ensure BlocksPerYear can be safely converted to int64
	if params.BlocksPerYear > uint64(stdmath.MaxInt64) {
		return sdk.Dec{}, fmt.Errorf("blocks per year (%d) exceeds maximum int64 value", params.BlocksPerYear)
	}

	// Check for division by zero and overflow
	if params.GoalBonded.IsZero() {
		return sdk.Dec{}, fmt.Errorf("goal bonded cannot be zero")
	}

	bondedRatioQuoGoalBonded := bondedRatio.Quo(params.GoalBonded)
	if bondedRatioQuoGoalBonded.GT(sdk.OneDec()) {
		bondedRatioQuoGoalBonded = sdk.OneDec()
	}

	inflationRateChangePerYear := sdk.OneDec().
		Sub(bondedRatioQuoGoalBonded).
		Mul(params.InflationRateChange)

	// Prevent division by zero
	if params.BlocksPerYear == 0 {
		return sdk.Dec{}, fmt.Errorf("blocks per year must be positive")
	}

	inflationRateChange := inflationRateChangePerYear.Quo(sdk.NewDec(int64(params.BlocksPerYear)))

	// adjust the new annual inflation for this next cycle
	inflation := m.Inflation.Add(inflationRateChange) // note inflationRateChange may be negative

	// Calculate min/max inflation bounds
	if totalStakingSupply.IsZero() {
		return sdk.ZeroDec(), nil
	}

	inflationMax := sdk.NewDecFromInt(params.MaxTokenPerYear).Quo(sdk.NewDecFromInt(totalStakingSupply))
	inflationMin := sdk.NewDecFromInt(params.MinTokenPerYear).Quo(sdk.NewDecFromInt(totalStakingSupply))

	// Validate the calculated inflation bounds
	if inflationMax.IsNil() || inflationMin.IsNil() {
		return sdk.Dec{}, fmt.Errorf("invalid inflation bounds calculated")
	}

	if inflation.GT(inflationMax) {
		inflation = inflationMax
	}
	if inflation.LT(inflationMin) {
		inflation = inflationMin
	}

	return inflation, nil
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply math.Int) (sdk.Dec, error) {
	if totalSupply.IsNil() {
		return sdk.Dec{}, fmt.Errorf("total supply cannot be nil")
	}
	if totalSupply.IsNegative() {
		return sdk.Dec{}, fmt.Errorf("total supply cannot be negative")
	}
	if totalSupply.GT(sdk.NewIntFromUint64(stdmath.MaxUint64)) {
		return sdk.Dec{}, fmt.Errorf("total supply too large")
	}
	if m.Inflation.IsNil() {
		return sdk.Dec{}, fmt.Errorf("inflation rate cannot be nil")
	}

	if totalSupply.IsZero() {
		return sdk.ZeroDec(), nil
	}

	return m.Inflation.MulInt(totalSupply), nil
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) (sdk.Coin, error) {
	// Validate inputs
	if m.AnnualProvisions.IsNegative() {
		return sdk.Coin{}, fmt.Errorf("annual provisions cannot be negative")
	}
	if params.MintDenom == "" {
		return sdk.Coin{}, fmt.Errorf("mint denom cannot be empty")
	}
	if err := sdk.ValidateDenom(params.MintDenom); err != nil {
		return sdk.Coin{}, fmt.Errorf("invalid mint denom: %w", err)
	}
	// Ensure BlocksPerYear can be safely converted to int64
	if params.BlocksPerYear > uint64(stdmath.MaxInt64) {
		return sdk.Coin{}, fmt.Errorf("blocks per year exceeds maximum int64 value")
	}
	if params.BlocksPerYear == 0 {
		return sdk.Coin{}, fmt.Errorf("blocks per year must be positive")
	}
	if m.AnnualProvisions.IsNil() {
		return sdk.Coin{}, fmt.Errorf("annual provisions cannot be nil")
	}

	provisionAmt := m.AnnualProvisions.QuoInt(sdk.NewInt(int64(params.BlocksPerYear)))
	if provisionAmt.IsNil() {
		return sdk.Coin{}, fmt.Errorf("invalid provision amount calculated")
	}

	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt()), nil
}

// CalculateInflationRate calculates the inflation rate for the current period
func (m Minter) CalculateInflationRate(params Params) (sdk.Dec, error) {
	// Validate inputs
	if params.InflationRateChange.IsNil() {
		return sdk.Dec{}, fmt.Errorf("inflation rate change cannot be nil")
	}
	if params.InflationRateChange.IsNegative() {
		return sdk.Dec{}, fmt.Errorf("inflation rate change cannot be negative")
	}
	// Ensure BlocksPerYear can be safely converted to int64
	if params.BlocksPerYear > uint64(stdmath.MaxInt64) {
		return sdk.Dec{}, fmt.Errorf("blocks per year exceeds maximum int64 value")
	}
	if params.BlocksPerYear == 0 {
		return sdk.Dec{}, fmt.Errorf("blocks per year cannot be zero")
	}
	inflationRateChangePerYear := params.InflationRateChange.Quo(sdk.NewDec(int64(params.BlocksPerYear)))
	if inflationRateChangePerYear.IsNil() {
		return sdk.Dec{}, fmt.Errorf("invalid inflation rate change calculated")
	}
	return m.Inflation.Add(inflationRateChangePerYear), nil
}
