package utils

import (
	"errors"
	"strings"

	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// OldBech32Prefix defines the Bech32 prefix used for EthAccounts
	OldBech32Prefix = "centauri"

	// OldBech32PrefixAccAddr defines the Bech32 prefix of an account's address
	OldBech32PrefixAccAddr = OldBech32Prefix
	// OldBech32PrefixAccPub defines the Bech32 prefix of an account's public key
	OldBech32PrefixAccPub = OldBech32Prefix + sdk.PrefixPublic
	// OldBech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	OldBech32PrefixValAddr = OldBech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// OldBech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	OldBech32PrefixValPub = OldBech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// OldBech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	OldBech32PrefixConsAddr = OldBech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// OldBech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	OldBech32PrefixConsPub = OldBech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func ConvertValAddr(valAddr string) string {
	parsedValAddr, err := ValAddressFromOldBech32(valAddr, OldBech32PrefixValAddr)
	if err != nil {
		return valAddr
	}
	return parsedValAddr.String()
}

func ConvertAccAddr(accAddr string) string {
	parsedAccAddr, err := AccAddressFromOldBech32(accAddr, OldBech32PrefixAccAddr)
	if err != nil {
		return accAddr
	}
	return parsedAccAddr.String()
}

func ConvertConsAddr(consAddr string) string {
	parsedConsAddr, err := ConsAddressFromOldBech32(consAddr, OldBech32PrefixConsAddr)
	if err != nil {
		return consAddr
	}
	return parsedConsAddr.String()
}

func IterateStoreByPrefix(
	ctx sdk.Context, storeKey storetypes.StoreKey, prefix []byte,
	fn func(value []byte) (newValue []byte),
) {
	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, prefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		newValue := fn(iterator.Value())
		store.Set(iterator.Key(), newValue)
	}
}

// AccAddressFromBech32 creates an AccAddress from a Bech32 string.
func AccAddressFromOldBech32(address, prefix string) (addr sdk.AccAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.AccAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.AccAddress(bz), nil
}

// ConsAddressFromBech32 creates a ConsAddress from a Bech32 string.
func ConsAddressFromOldBech32(address, prefix string) (addr sdk.ConsAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ConsAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.ConsAddress(bz), nil
}

// ValAddressFromBech32 creates a ValAddress from a Bech32 string.
func ValAddressFromOldBech32(address, prefix string) (addr sdk.ValAddress, err error) {
	if len(strings.TrimSpace(address)) == 0 {
		return sdk.ValAddress{}, errors.New("empty address string is not allowed")
	}

	bz, err := sdk.GetFromBech32(address, prefix)
	if err != nil {
		return nil, err
	}

	err = sdk.VerifyAddressFormat(bz)
	if err != nil {
		return nil, err
	}

	return sdk.ValAddress(bz), nil
}
