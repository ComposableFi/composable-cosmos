package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	mintkeeper "github.com/notional-labs/composable/v6/x/mint/keeper"
)

type Keeper struct {
	stakingkeeper.Keeper
	keys       storetypes.StoreKey
	cdc        codec.BinaryCodec
	mintkeeper mintkeeper.Keeper
	acck       accountkeeper.AccountKeeper
	authority  string
}

// func NewBaseKeeper(
// 	cdc codec.BinaryCodec,
// 	key storetypes.StoreKey,
// 	ak types.AccountKeeper,
// 	acck accountkeeper.AccountKeeper,
// 	bk bankkeeper.Keeper,
// 	authority string,
// ) Keeper {
// 	keeper := Keeper{
// 		Keeper: *stakingkeeper.NewKeeper(cdc, key, ak, bk, authority),
// 		acck:   acck,
// 	}
// 	return keeper
// }

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	keys storetypes.StoreKey,
	staking stakingkeeper.Keeper,
	acck accountkeeper.AccountKeeper,
	mintkeeper mintkeeper.Keeper,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper:     staking,
		keys:       keys,
		acck:       acck,
		authority:  authority,
		mintkeeper: mintkeeper,
	}
	return keeper
}

// func (k *Keeper) RegisterKeepers(akk banktypes.StakingKeeper) {
// 	k.acck = sk
// }

// func (k Keeper) StoreDelegation(ctx sdk.Context, delegation types.Delegation) {
// 	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

// 	store := ctx.KVStore(k.storeKey)
// 	b := types.MustMarshalDelegation(k.cdc, delegation)
// 	store.Set(customstakingtypes.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()), b)
// }
