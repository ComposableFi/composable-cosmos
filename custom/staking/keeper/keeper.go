package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	mintkeeper "github.com/notional-labs/composable/v6/x/mint/keeper"
)

type Keeper struct {
	stakingkeeper.Keeper
	cdc        codec.BinaryCodec
	acck       accountkeeper.AccountKeeper
	mintkeeper *mintkeeper.Keeper
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

func NewKeeper(
	cdc codec.BinaryCodec,
	staking stakingkeeper.Keeper,
	acck accountkeeper.AccountKeeper,
	mintkeeper *mintkeeper.Keeper,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper:     staking,
		acck:       acck,
		authority:  authority,
		mintkeeper: mintkeeper,
		cdc:        cdc,
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
