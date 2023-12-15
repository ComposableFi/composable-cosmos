package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
)

// MinterKey is the key to use for the keeper store.
var (
	DelegationKey = []byte{0x02} // key for a delegation
)

const (
	// module name
	ModuleName = "mint"

	// StoreKey is the default store key for mint
	StoreKey = ModuleName

	// QuerierRoute is the querier route for the minting store.
	QuerierRoute = StoreKey

	// Query endpoints supported by the minting querier
	QueryParameters       = "parameters"
	QueryInflation        = "inflation"
	QueryAnnualProvisions = "annual_provisions"
)

// GetDelegationKey creates the key for delegator bond with validator
// VALUE: staking/Delegation
func GetDelegationKey(delAddr sdk.AccAddress, valAddr sdk.ValAddress) []byte {
	return append(GetDelegationsKey(delAddr), address.MustLengthPrefix(valAddr)...)
}

func GetDelegationsKey(delAddr sdk.AccAddress) []byte {
	return append(DelegationKey, address.MustLengthPrefix(delAddr)...)
}
