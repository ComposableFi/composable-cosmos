package utils

import (
	"errors"
	"strings"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
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

	// OldBech32Prefix defines the Bech32 prefix used for EthAccounts
	NewBech32Prefix = "pica"

	// NewBech32PrefixAccAddr defines the Bech32 prefix of an account's address
	NewBech32PrefixAccAddr = NewBech32Prefix
	// NewBech32PrefixAccPub defines the Bech32 prefix of an account's public key
	NewBech32PrefixAccPub = NewBech32Prefix + sdk.PrefixPublic
	// NewBech32PrefixValAddr defines the Bech32 prefix of a validator's operator address
	NewBech32PrefixValAddr = NewBech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator
	// NewBech32PrefixValPub defines the Bech32 prefix of a validator's operator public key
	NewBech32PrefixValPub = NewBech32Prefix + sdk.PrefixValidator + sdk.PrefixOperator + sdk.PrefixPublic
	// NewBech32PrefixConsAddr defines the Bech32 prefix of a consensus node address
	NewBech32PrefixConsAddr = NewBech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus
	// NewBech32PrefixConsPub defines the Bech32 prefix of a consensus node public key
	NewBech32PrefixConsPub = NewBech32Prefix + sdk.PrefixValidator + sdk.PrefixConsensus + sdk.PrefixPublic
)

func ConvertValAddr(valAddr string) string {
	parsedValAddr, err := ValAddressFromOldBech32(valAddr, OldBech32PrefixValAddr)
	_, bz, _ := bech32.DecodeAndConvert(parsedValAddr.String())
	bech32Addr, _ := bech32.ConvertAndEncode(NewBech32PrefixValAddr, bz)
	if err != nil {
		return valAddr
	}
	return bech32Addr
}

func ConvertAccAddr(accAddr string) string {
	parsedAccAddr, err := AccAddressFromOldBech32(accAddr, OldBech32PrefixAccAddr)
	_, bz, _ := bech32.DecodeAndConvert(parsedAccAddr.String())
	bech32Addr, _ := bech32.ConvertAndEncode(NewBech32PrefixAccAddr, bz)
	if err != nil {
		panic(err)
	}
	return bech32Addr
}

// Input is type string -> need safe convert
func SafeConvertAddress(accAddr string) string {
	if len(accAddr) == 0 {
		return ""
	}

	parsedAccAddr, err := AccAddressFromOldBech32(accAddr, OldBech32PrefixAccAddr)
	if err != nil {
		return accAddr
	}
	_, bz, err := bech32.DecodeAndConvert(parsedAccAddr.String())
	if err != nil {
		return accAddr
	}
	bech32Addr, err := bech32.ConvertAndEncode(NewBech32PrefixAccAddr, bz)
	if err != nil {
		return accAddr
	}

	return bech32Addr
}

func ConvertConsAddr(consAddr string) string {
	parsedConsAddr, err := ConsAddressFromOldBech32(consAddr, OldBech32PrefixConsAddr)
	_, bz, _ := bech32.DecodeAndConvert(parsedConsAddr.String())
	bech32Addr, _ := bech32.ConvertAndEncode(NewBech32PrefixConsAddr, bz)
	if err != nil {
		return consAddr
	}
	return bech32Addr
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
