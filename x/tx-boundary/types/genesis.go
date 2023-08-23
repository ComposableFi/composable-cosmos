package types

import fmt "fmt"

var (
	DefaultDelegateBoundary = Boundary{
		TxLimit:             5,
		BlocksPerGeneration: 5,
	}
	DefaultRedelegateBoundary = Boundary{
		TxLimit:             5,
		BlocksPerGeneration: 5,
	}
)

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		DelegateBoundary:   DefaultDelegateBoundary,
		RedelegateBoundary: DefaultRedelegateBoundary,
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func ValidateGenesis(data GenesisState) error {
	if data.DelegateBoundary.BlocksPerGeneration <= 0 || data.RedelegateBoundary.BlocksPerGeneration <= 0 {
		return fmt.Errorf("BlocksPerGeneration must greater than 0")
	}
	return nil
}
