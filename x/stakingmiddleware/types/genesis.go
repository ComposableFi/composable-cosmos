package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(m math.Int, params Params) *GenesisState {
	return &GenesisState{
		LastTotalPower: m,
		Params:         params,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		LastTotalPower: sdk.ZeroInt(),
		Params:         Params{BlocksPerEpoch: 10},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return nil
}
