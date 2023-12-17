package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// MinterKey is the key to use for the keeper store.
var (
	DelegationKey = []byte{0x01} // key for a delegation
)

const (
	// module name
	ModuleName = "stakingmiddleware"

	// StoreKey is the default store key for mint

	StoreKey = "customstmiddleware" // not using the module name because of collisions with key "ibc"
)

// GetDelegationKey creates the key for delegator bond with validator
// VALUE: staking/Delegation
func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, address.MustLengthPrefix(delAddr)...)
}
