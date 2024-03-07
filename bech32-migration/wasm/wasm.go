package wasm

import (
	"github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/notional-labs/composable/v6/bech32-migration/utils"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func MigrateAddressBech32(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	migrateCodeInfo(ctx, storeKey, cdc)
	migrateContractInfo(ctx, storeKey, cdc)
}

func migrateCodeInfo(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	// Code id
	ctx.Logger().Debug("Migrationg of address bech32 for wasm module Code Info begin")
	prefixStore := prefix.NewStore(ctx.KVStore(storeKey), types.CodeKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)
	defer iter.Close()

	totalMigratedCodeId := uint64(0)
	for ; iter.Valid(); iter.Next() {
		// get code info value
		var c types.CodeInfo
		cdc.MustUnmarshal(iter.Value(), &c)

		// Update info
		c.Creator = utils.SafeConvertAddress(c.Creator)
		c.InstantiateConfig.Address = utils.SafeConvertAddress(c.InstantiateConfig.Address)
		for i := range c.InstantiateConfig.Addresses {
			c.InstantiateConfig.Addresses[i] = utils.SafeConvertAddress(c.InstantiateConfig.Addresses[i])
		}

		// save updated code info
		prefixStore.Set(iter.Key(), cdc.MustMarshal(&c))

		totalMigratedCodeId++
	}

	// contract info prefix store
	ctx.Logger().Debug(
		"Migration of address bech32 for wasm module code info done",
		"total_migrated_code_id", totalMigratedCodeId,
	)
}

func migrateContractInfo(ctx sdk.Context, storeKey storetypes.StoreKey, cdc codec.BinaryCodec) {
	ctx.Logger().Debug("Migrating of addresses bech32 for wasm module Contract info begin")
	// contract info prefix store
	prefixStore := prefix.NewStore(ctx.KVStore(storeKey), types.ContractKeyPrefix)
	iter := prefixStore.Iterator(nil, nil)

	defer iter.Close()

	totalMigratedContractAddresses := uint64(0)
	for ; iter.Valid(); iter.Next() {
		// get code info value
		var c types.ContractInfo
		cdc.MustUnmarshal(iter.Value(), &c)

		// Update info
		c.Creator = utils.SafeConvertAddress(c.Creator)
		c.Admin = utils.SafeConvertAddress(c.Admin)
		// save updated code info
		prefixStore.Set(iter.Key(), cdc.MustMarshal(&c))

		totalMigratedContractAddresses++
	}

	ctx.Logger().Debug(
		"Migrating of addresses bech32 for wasm module Contract info done",
		"total_migrated_contract_addresses", totalMigratedContractAddresses,
	)
}
