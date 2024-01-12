package types

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params, rewardDenom RewardDenom) *GenesisState {
	return &GenesisState{
		Params:      params,
		RewardDenom: rewardDenom,
	}
}

// DefaultGenesisState creates a default GenesisState object
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params: Params{BlocksPerEpoch: 10, AllowUnbondAfterEpochProgressBlockNumber: 0},
		// need to change to ppica for mainnet
		RewardDenom: RewardDenom{Denom: "stake"},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	return nil
}
