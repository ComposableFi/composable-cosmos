package keeper

import (
	"fmt"

	"github.com/cometbft/cometbft/libs/log"
	"github.com/notional-labs/composable/v6/x/stakingmiddleware/types"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper of the staking middleware store
type Keeper struct {
	cdc           codec.BinaryCodec
	storeKey      storetypes.StoreKey
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates a new middleware Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:           cdc,
		storeKey:      key,
		accountKeeper: ak,
		bankKeeper:    bk,
		authority:     authority,
	}
}

// GetAuthority returns the x/stakingmiddleware module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// SetParams sets the x/stakingmiddleware module parameters.
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	if p.BlocksPerEpoch < 5 {
		return fmt.Errorf(
			"BlocksPerEpoch must be greater than or equal to 5",
		)
	}
	if p.AllowUnbondAfterEpochProgressBlockNumber > p.BlocksPerEpoch {
		return fmt.Errorf(
			"AllowUnbondAfterEpochProgressBlockNumber must be less than or equal to BlocksPerEpoch",
		)
	}

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&p)
	store.Set(types.ParamsKey, bz)

	return nil
}

// SetParams sets the x/stakingmiddleware module parameters.
func (k Keeper) SetDenom(ctx sdk.Context, rd types.RewardDenom) error {

	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&rd)
	store.Set(types.RewardDenomKey, bz)

	return nil
}

// GetParams returns the current x/stakingmiddleware module parameters.
func (k Keeper) GetRewardDenom(ctx sdk.Context) (rd types.RewardDenom) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.RewardDenomKey)
	if bz == nil {
		return rd
	}

	k.cdc.MustUnmarshal(bz, &rd)
	return rd
}

// GetParams returns the current x/stakingmiddleware module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

func (k Keeper) GetModuleAccountAccAddress(ctx sdk.Context) sdk.AccAddress {
	moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.RewardModuleName)
	return moduleAccount.GetAddress()
}
