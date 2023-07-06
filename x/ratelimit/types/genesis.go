package types

// DefaultGenesis returns the default Capability genesis state
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:     DefaultParams(),
		RateLimits: []RateLimit{},
		Epochs:     []EpochInfo{NewGenesisEpochInfo(HOUR_EPOCH, EpochHourPeriod)},
	}
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func ValidateGenesis(data GenesisState) error {
	return data.Params.Validate()
}
