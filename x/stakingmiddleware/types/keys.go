package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// MinterKey is the key to use for the keeper store.
var (
	DelegateKey                  = []byte{0x01} // key for a delegation
	BeginRedelegateKey           = []byte{0x02} // key for a delegation
	UndelegateKey                = []byte{0x03} // key for a delegation
	CancelUnbondingDelegationKey = []byte{0x04} // key for a delegation
	ParamsKey                    = []byte{0x05} // key for global staking middleware params in the keeper store
)

const (
	// module name
	ModuleName = "stakingmiddleware"

	// StoreKey is the default store key for mint

	StoreKey = "customstmiddleware" // not using the module name because of collisions with key "ibc"
)

// GetDelegateKey creates the key for delegator bond with validator
// VALUE: staking/Delegation
func GetDelegateKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(DelegationKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func DelegationKey(delAddr sdk.AccAddress) []byte {
	return append(DelegateKey, address.MustLengthPrefix(delAddr)...)
}

func GetBeginRedelegateKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(BeginRedelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func BeginRedelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(BeginRedelegateKey, address.MustLengthPrefix(delAddr)...)
}

func GetUndelegateKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(UndelegateionKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func UndelegateionKey(delAddr sdk.AccAddress) []byte {
	return append(UndelegateKey, address.MustLengthPrefix(delAddr)...)
}

func GetCancelUnbondingDelegateKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(CancelUnbondingDelegateKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func CancelUnbondingDelegateKey(delAddr sdk.AccAddress) []byte {
	return append(CancelUnbondingDelegationKey, address.MustLengthPrefix(delAddr)...)
}
