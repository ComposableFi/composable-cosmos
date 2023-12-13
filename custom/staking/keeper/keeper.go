package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	customstakingtypes "github.com/notional-labs/composable/v6/custom/staking/types"
)

type Keeper struct {
	stakingkeeper.Keeper
	cdc       codec.BinaryCodec
	storeKey  storetypes.StoreKey
	acck      accountkeeper.AccountKeeper
	authority string
}

var _ stakingkeeper.Keeper = stakingkeeper.Keeper{} //???

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

func NewBaseKeeper2(
	cdc codec.BinaryCodec,
	keys storetypes.StoreKey,
	staking stakingkeeper.Keeper,
	acck accountkeeper.AccountKeeper,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper:    staking,
		acck:      acck,
		authority: authority,
	}
	return keeper
}

// func (k *Keeper) RegisterKeepers(akk banktypes.StakingKeeper) {
// 	k.acck = sk
// }

func (k Keeper) StoreDelegation(ctx sdk.Context, delegation types.Delegation) {
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	b := types.MustMarshalDelegation(k.cdc, delegation)
	store.Set(customstakingtypes.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr()), b)
}
