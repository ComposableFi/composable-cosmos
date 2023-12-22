package keeper

import (
	abcicometbft "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	"github.com/cosmos/cosmos-sdk/x/staking/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	mintkeeper "github.com/notional-labs/composable/v6/x/mint/keeper"
	stakingmiddleware "github.com/notional-labs/composable/v6/x/stakingmiddleware/keeper"
)

type Keeper struct {
	stakingkeeper.Keeper
	cdc               codec.BinaryCodec
	acck              accountkeeper.AccountKeeper
	mintkeeper        *mintkeeper.Keeper
	Stakingmiddleware *stakingmiddleware.Keeper
	authority         string
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

func (k Keeper) BlockValidatorUpdates(ctx sdk.Context, hight int64) []abcicometbft.ValidatorUpdate {
	// Calculate validator set changes.
	//
	// NOTE: ApplyAndReturnValidatorSetUpdates has to come before
	// UnbondAllMatureValidatorQueue.
	// This fixes a bug when the unbonding period is instant (is the case in
	// some of the tests). The test expected the validator to be completely
	// unbonded after the Endblocker (go from Bonded -> Unbonding during
	// ApplyAndReturnValidatorSetUpdates and then Unbonding -> Unbonded during
	// UnbondAllMatureValidatorQueue).
	println("BlockValidatorUpdates Custom Staking Module")
	params := k.Stakingmiddleware.GetParams(ctx)
	println("BlocksPerEpoch: ", params.BlocksPerEpoch)
	should_execute_batch := (hight % int64(params.BlocksPerEpoch)) == 0
	var validatorUpdates []abcicometbft.ValidatorUpdate
	if should_execute_batch {
		println("Should Execute Batch: ", hight)
		v, err := k.ApplyAndReturnValidatorSetUpdates(ctx)
		if err != nil {
			panic(err)
		}
		validatorUpdates = v
	}

	// unbond all mature validators from the unbonding queue
	k.UnbondAllMatureValidators(ctx)

	// Remove all mature unbonding delegations from the ubd queue.
	matureUnbonds := k.DequeueAllMatureUBDQueue(ctx, ctx.BlockHeader().Time)
	for _, dvPair := range matureUnbonds {
		addr, err := sdk.ValAddressFromBech32(dvPair.ValidatorAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvPair.DelegatorAddress)

		balances, err := k.CompleteUnbonding(ctx, delegatorAddress, addr)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteUnbonding,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
			),
		)
	}

	// Remove all mature redelegations from the red queue.
	matureRedelegations := k.DequeueAllMatureRedelegationQueue(ctx, ctx.BlockHeader().Time)
	for _, dvvTriplet := range matureRedelegations {
		valSrcAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorSrcAddress)
		if err != nil {
			panic(err)
		}
		valDstAddr, err := sdk.ValAddressFromBech32(dvvTriplet.ValidatorDstAddress)
		if err != nil {
			panic(err)
		}
		delegatorAddress := sdk.MustAccAddressFromBech32(dvvTriplet.DelegatorAddress)

		balances, err := k.CompleteRedelegation(
			ctx,
			delegatorAddress,
			valSrcAddr,
			valDstAddr,
		)
		if err != nil {
			continue
		}

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeCompleteRedelegation,
				sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
				sdk.NewAttribute(types.AttributeKeyDelegator, dvvTriplet.DelegatorAddress),
				sdk.NewAttribute(types.AttributeKeySrcValidator, dvvTriplet.ValidatorSrcAddress),
				sdk.NewAttribute(types.AttributeKeyDstValidator, dvvTriplet.ValidatorDstAddress),
			),
		)
	}

	return validatorUpdates
}

func NewKeeper(
	cdc codec.BinaryCodec,
	staking stakingkeeper.Keeper,
	acck accountkeeper.AccountKeeper,
	mintkeeper *mintkeeper.Keeper,
	stakingmiddleware *stakingmiddleware.Keeper,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper:            staking,
		acck:              acck,
		authority:         authority,
		mintkeeper:        mintkeeper,
		Stakingmiddleware: stakingmiddleware,
		cdc:               cdc,
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
