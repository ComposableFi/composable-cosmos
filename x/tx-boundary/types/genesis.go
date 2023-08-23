package types

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
	return nil
}
