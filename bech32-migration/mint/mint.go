package mint

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/notional-labs/composable/v6/bech32-migration/utils"
	"github.com/notional-labs/composable/v6/x/mint/types"
)

func MigrateAddressBech32(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	ctx.Logger().Info("Migration of address bech32 for mint module begin")
	interchainAccountCount := uint64(0)

	store := ctx.KVStore(storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.AllowedAddressKey)

	for ; iterator.Valid(); iterator.Next() {
		interchainAccountCount++
		trimedAddr := strings.Replace(string(iterator.Key()), "\x01", "", 1)
		newPrefixAddr := utils.ConvertAccAddr(trimedAddr)
		key := types.GetAllowedAddressStoreKey(newPrefixAddr)
		store.Set(key, []byte{1})
	}

	ctx.Logger().Info(
		"Migration of address bech32 for mint module done",
		"key_changed_count", interchainAccountCount,
	)
}
