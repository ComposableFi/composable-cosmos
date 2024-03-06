package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/notional-labs/composable/v6/x/keyrotation/types"
)

type Keeper struct {
	storeKey storetypes.StoreKey
	cdc      codec.BinaryCodec
	sk       types.StakingKeeper

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	sk types.StakingKeeper,
	authority string,
) Keeper {
	return Keeper{
		cdc:       cdc,
		storeKey:  key,
		sk:        sk,
		authority: authority,
	}
}

func (k Keeper) handleMsgRotateConsPubKey(ctx sdk.Context, valAddress sdk.ValAddress, pubKey cryptotypes.PubKey) error {
	var (
		validator stakingtypes.Validator
		found     bool
	)
	// check to see if the validator not exist
	if validator, found = k.sk.GetValidator(ctx, valAddress); !found {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "validator not exists")
	}

	// check if pubkey is exists
	if _, found := k.sk.GetValidatorByConsAddr(ctx, sdk.GetConsAddress(pubKey)); found {
		return stakingtypes.ErrValidatorPubKeyExists
	}

	// check if new pubkey is allowed
	cp := ctx.ConsensusParams()
	if cp != nil && cp.Validator != nil {
		pkType := pubKey.Type()
		hasKeyType := false
		for _, keyType := range cp.Validator.PubKeyTypes {
			if pkType == keyType {
				hasKeyType = true
				break
			}
		}
		if !hasKeyType {
			return errorsmod.Wrapf(
				stakingtypes.ErrValidatorPubKeyTypeNotSupported,
				"got: %s, expected: %s", pubKey.Type(), cp.Validator.PubKeyTypes,
			)
		}
	}

	// wrap pubkey to types any
	newPkAny, err := codectypes.NewAnyWithValue(pubKey)
	if err != nil {
		return err
	}

	oldPkAny := validator.ConsensusPubkey

	// replace pubkey
	validator.ConsensusPubkey = newPkAny

	// NOTE: staking module do not support DeleteValidatorByConsAddr method for Validator types
	// we need to record all the rotated pubkey in keyrotation store so when we can delete later
	// by upgrade handler

	// set validator
	k.sk.SetValidator(ctx, validator)
	k.sk.SetValidatorByConsAddr(ctx, validator)

	// Set rotation history
	consPubKeyRotationHistory := types.ConsPubKeyRotationHistory{
		OperatorAddress: valAddress.String(),
		OldKey:          oldPkAny,
		NewKey:          newPkAny,
		BlockHeight:     uint64(ctx.BlockHeight()),
	}

	k.SetKeyRotationHistory(ctx, consPubKeyRotationHistory)

	return nil
}
