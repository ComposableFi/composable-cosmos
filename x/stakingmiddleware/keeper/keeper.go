package keeper

import (
	"github.com/cometbft/cometbft/libs/log"
	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkmath "cosmossdk.io/math"
)

// Keeper of the staking middleware store
type Keeper struct {
	cdc      codec.BinaryCodec
	storeKey storetypes.StoreKey
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates a new middleware Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	authority string,
) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		authority: authority,
	}
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// // SetParams sets the x/mint module parameters.
// func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
// 	if err := p.Validate(); err != nil {
// 		return err
// 	}

// 	store := ctx.KVStore(k.storeKey)
// 	bz := k.cdc.MustMarshal(&p)
// 	store.Set(types.ParamsKey, bz)

// 	return nil
// }

// // GetParams returns the current x/mint module parameters.
// func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
// 	store := ctx.KVStore(k.storeKey)
// 	bz := store.Get(types.ParamsKey)
// 	if bz == nil {
// 		return p
// 	}

// 	k.cdc.MustUnmarshal(bz, &p)
// 	return p
// }

// func (k Keeper) StoreDelegation(ctx sdk.Context, delegation stakingtypes.Delegation) {
// 	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)
// 	log := k.Logger(ctx)
// 	log.Info("StoreDelegation", "delegatorAddress", delegatorAddress, "validatorAddress", delegation.GetValidatorAddr())
// 	store := ctx.KVStore(k.storeKey)
// 	b := stakingtypes.MustMarshalDelegation(k.cdc, delegation)
// 	kkk := types.GetDelegationKey(delegatorAddress, delegation.GetValidatorAddr())
// 	// log.Info()
// 	store.Set(kkk, b)
// }

// SetLastTotalPower Set the last total validator power.
func (k Keeper) SetLastTotalPower(ctx sdk.Context, power sdkmath.Int) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&sdk.IntProto{Int: power})
	store.Set(types.DelegationKey, bz)
}

func (k Keeper) GetLastTotalPower(ctx sdk.Context) sdkmath.Int {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.DelegationKey)

	if bz == nil {
		return sdkmath.ZeroInt()
	}

	ip := sdk.IntProto{}
	k.cdc.MustUnmarshal(bz, &ip)

	return ip.Int
}

func (k Keeper) SetDelegation(ctx sdk.Context, sourceDelegatorAddress, validatorAddress, denom string, amount sdkmath.Int) {
	delegation := types.Delegation{DelegatorAddress: sourceDelegatorAddress, ValidatorAddress: validatorAddress, Amount: sdk.NewCoin(denom, amount)}
	delegatorAddress := sdk.MustAccAddressFromBech32(delegation.DelegatorAddress)

	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&delegation)
	store.Set(types.GetDelegationKey(delegatorAddress, GetValidatorAddr(delegation)), b)
}

func (k Keeper) IterateDelegations(ctx sdk.Context, fn func(index int64, ubd types.Delegation) (stop bool)) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.DelegationKey)
	defer iterator.Close()

	for i := int64(0); iterator.Valid(); iterator.Next() {
		ubd := MustUnmarshalUBD(k.cdc, iterator.Value())
		if stop := fn(i, ubd); stop {
			break
		}
		i++
	}
}

// DequeueAllMatureUBDQueue returns a concatenated list of all the timeslices inclusively previous to
// currTime, and deletes the timeslices from the queue.
func (k Keeper) DequeueAllDelegation(ctx sdk.Context) (delegations []types.Delegation) {
	store := ctx.KVStore(k.storeKey)

	// gets an iterator for all timeslices from time 0 until the current Blockheader time
	delegationIterator := sdk.KVStorePrefixIterator(store, types.DelegationKey)
	defer delegationIterator.Close()

	for ; delegationIterator.Valid(); delegationIterator.Next() {
		delegation := types.Delegation{}
		value := delegationIterator.Value()
		k.cdc.MustUnmarshal(value, &delegation)

		delegations = append(delegations, delegation)

		store.Delete(delegationIterator.Key())
	}

	return delegations
}

func GetValidatorAddr(d types.Delegation) sdk.ValAddress {
	addr, err := sdk.ValAddressFromBech32(d.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return addr
}

func UnmarshalBD(cdc codec.BinaryCodec, value []byte) (ubd types.Delegation, err error) {
	err = cdc.Unmarshal(value, &ubd)
	return ubd, err
}

func MustUnmarshalUBD(cdc codec.BinaryCodec, value []byte) types.Delegation {
	ubd, err := UnmarshalBD(cdc, value)
	if err != nil {
		panic(err)
	}

	return ubd
}
