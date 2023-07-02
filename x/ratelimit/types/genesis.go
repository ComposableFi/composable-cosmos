package types

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		RateLimits: []RateLimit{},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
