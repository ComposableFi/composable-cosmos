package types

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func ValidateGenesis(data GenesisState) error {
	return nil
}
