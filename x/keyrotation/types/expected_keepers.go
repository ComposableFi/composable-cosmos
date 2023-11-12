package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// StakingKeeper expected staking keeper
type StakingKeeper interface {
	SetValidator(sdk.Context, stakingtypes.Validator)
	SetValidatorByConsAddr(sdk.Context, stakingtypes.Validator) error
	GetValidator(sdk.Context, sdk.ValAddress) (stakingtypes.Validator, bool)
	GetValidatorByConsAddr(sdk.Context, sdk.ConsAddress) (stakingtypes.Validator, bool)
}
