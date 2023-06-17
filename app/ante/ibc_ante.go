package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/authz"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
)

var (
	allowedRelayAddress = map[string]bool{
		"centauri1eqv3xl0vk0md74qukfghfff4z3axsp29rr9c85": true,
		"centauri1av6x9sll0yx4anske424jtgxejnrgqv6j6tjjt": true,
	}
)

type IBCPermissionDecorator struct {
	cdc codec.BinaryCodec
}

func NewIBCPermissionDecorator(cdc codec.BinaryCodec) IBCPermissionDecorator {
	return IBCPermissionDecorator{
		cdc: cdc,
	}
}

func (g IBCPermissionDecorator) AnteHandle(
	ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	// run checks only on CheckTx or simulate
	if !ctx.IsCheckTx() || simulate {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	if err = g.ValidateIBCUpdateClientMsg(ctx, msgs); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

// ValidateIBCUpdateClientMsg validate
func (g IBCPermissionDecorator) ValidateIBCUpdateClientMsg(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, m := range msgs {
		if msg, ok := m.(*authz.MsgExec); ok {
			if err := g.validAuthz(msg); err != nil {
				return err
			}
			continue
		}

		// validate normal msgs
		if err := g.validMsg(m); err != nil {
			return err
		}
	}
	return nil
}

func (g IBCPermissionDecorator) validMsg(m sdk.Msg) error {
	if msg, ok := m.(*clienttypes.MsgUpdateClient); ok {
		if !allowedRelayAddress[msg.Signer] {
			return fmt.Errorf("permission denied, address %s don't have relay permission", msg.Signer)
		}
		// prevent messages with insufficient initial deposit amount
	}

	return nil
}

func (g IBCPermissionDecorator) validAuthz(execMsg *authz.MsgExec) error {
	for _, v := range execMsg.Msgs {
		var innerMsg sdk.Msg
		if err := g.cdc.UnpackAny(v, &innerMsg); err != nil {
			return errorsmod.Wrap(err, "cannot unmarshal authz exec msgs")
		}
		if err := g.validMsg(innerMsg); err != nil {
			return err
		}
	}

	return nil
}
