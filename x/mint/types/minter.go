package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewMinter returns a new Minter object with the given inflation and annual
// provisions values.
func NewMinter(inflation, annualProvisions sdkmath.LegacyDec) Minter {
	return Minter{
		Inflation:        inflation,
		AnnualProvisions: annualProvisions,
	}
}

// InitialMinter returns an initial Minter object with a given inflation value.
func InitialMinter(inflation sdkmath.LegacyDec) Minter {
	return NewMinter(
		inflation,
		sdkmath.LegacyNewDec(0),
	)
}

// DefaultInitialMinter returns a default initial Minter object for a new chain
// which uses an inflation rate of 13% per year.
func DefaultInitialMinter() Minter {
	return InitialMinter(
		// Create a new Dec from integer with decimal place at prec
		// CONTRACT: prec <= Precision
		sdkmath.LegacyNewDecWithPrec(InflationRate, Precision),
	)
}

// validate minter
func ValidateMinter(minter Minter) error {
	if minter.Inflation.IsNegative() {
		return fmt.Errorf("mint parameter Inflation should be positive, is %s",
			minter.Inflation.String())
	}
	return nil
}

// NextInflationRate returns the new inflation rate for the next hour.
func (m Minter) NextInflationRate(params Params, bondedRatio sdkmath.LegacyDec, totalStakingSupply sdkmath.Int) sdkmath.LegacyDec {
	totalStakingSupplyDec := sdkmath.LegacyNewDecFromInt(totalStakingSupply)
	if totalStakingSupplyDec.LT(math.LegacySmallestDec()) {
		return m.Inflation // assert if totalStakingSupplyDec = 0
	}

	// The target annual inflation rate is recalculated for each previsions cycle. The
	// inflation is also subject to a rate change (positive or negative) depending on
	// the distance from the desired ratio (67%). The maximum rate change possible is
	// defined to be 13% per year, however the annual inflation is capped as between
	// 7% and 20%.

	// (1 - bondedRatio/GoalBonded) * InflationRateChange
	inflationRateChangePerYear := sdkmath.LegacyOneDec().
		Sub(bondedRatio.Quo(params.GoalBonded)).
		Mul(params.InflationRateChange)
	inflationRateChange := inflationRateChangePerYear.Quo(sdkmath.LegacyNewDec(int64(params.BlocksPerYear)))

	// adjust the new annual inflation for this next cycle
	inflation := m.Inflation.Add(inflationRateChange) // note inflationRateChange may be negative

	inflationMax := sdkmath.LegacyNewDecFromInt(params.MaxTokenPerYear).Quo(totalStakingSupplyDec)
	inflationMin := sdkmath.LegacyNewDecFromInt(params.MinTokenPerYear).Quo(totalStakingSupplyDec)

	if inflation.GT(inflationMax) {
		inflation = inflationMax
	}
	if inflation.LT(inflationMin) {
		inflation = inflationMin
	}

	return inflation
}

// NextAnnualProvisions returns the annual provisions based on current total
// supply and inflation rate.
func (m Minter) NextAnnualProvisions(_ Params, totalSupply math.Int) sdkmath.LegacyDec {
	return m.Inflation.MulInt(totalSupply)
}

// BlockProvision returns the provisions for a block based on the annual
// provisions rate.
func (m Minter) BlockProvision(params Params) sdk.Coin {
	provisionAmt := m.AnnualProvisions.QuoInt(sdkmath.NewInt(int64(params.BlocksPerYear)))
	return sdk.NewCoin(params.MintDenom, provisionAmt.TruncateInt())
}
