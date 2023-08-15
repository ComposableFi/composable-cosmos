package types

import (
	fmt "fmt"

	"cosmossdk.io/math"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	DefaultMinRateLimitAmount = math.NewIntFromUint64(10_000_000_000)
)

// Param keys store keys
var (
	KeyMinRateLimitAmount = []byte("minratelimitamount")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(minRateLimitAmount math.Int) Params {
	return Params{MinRateLimitAmount: minRateLimitAmount}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return NewParams(DefaultMinRateLimitAmount)
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyMinRateLimitAmount, &p.MinRateLimitAmount, validateMinRateLimitAmount),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	return nil
}

func validateMinRateLimitAmount(i interface{}) error {
	_, ok := i.(math.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type string: %T", i)
	}

	return nil
}
