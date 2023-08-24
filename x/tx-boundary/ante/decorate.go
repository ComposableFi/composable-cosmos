package ante

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	txBoundaryKeeper "github.com/notional-labs/centauri/v4/x/tx-boundary/keeper"
)

type StakingPermissionDecorator struct {
	cdc        codec.BinaryCodec
	txBoundary txBoundaryKeeper.Keeper
}

func NewStakingPermissionDecorator(cdc codec.BinaryCodec, keeper txBoundaryKeeper.Keeper) StakingPermissionDecorator {
	return StakingPermissionDecorator{
		cdc:        cdc,
		txBoundary: keeper,
	}
}

func (g StakingPermissionDecorator) AnteHandle(
	ctx sdk.Context, tx sdk.Tx,
	simulate bool, next sdk.AnteHandler,
) (newCtx sdk.Context, err error) {
	// run checks only on CheckTx or simulate
	if simulate {
		return next(ctx, tx, simulate)
	}

	msgs := tx.GetMsgs()
	if err = g.ValidateStakingMsgs(ctx, msgs); err != nil {
		return ctx, err
	}

	return next(ctx, tx, simulate)
}

// ValidateStakingMsg validate
func (g StakingPermissionDecorator) ValidateStakingMsgs(ctx sdk.Context, msgs []sdk.Msg) error {
	for _, m := range msgs {
		g.ValidateStakingMsg(ctx, m)
	}
	return nil
}

func (g StakingPermissionDecorator) ValidateStakingMsg(ctx sdk.Context, msg sdk.Msg) error {
	switch msg := msg.(type) {

	case *stakingtypes.MsgDelegate:
		if err := g.validDelegateMsg(ctx, msg); err != nil {
			return err
		}
	case *stakingtypes.MsgBeginRedelegate:
		if err := g.validRedelegateMsg(ctx, msg); err != nil {
			return err
		}
	case *authz.MsgExec:
		if err := g.validAuthz(ctx, msg); err != nil {
			return err
		}
	default:
		return nil
	}
	return nil
}

func (g StakingPermissionDecorator) validDelegateMsg(ctx sdk.Context, msg *stakingtypes.MsgDelegate) error {
	boundary := g.txBoundary.GetDelegateBoundary(ctx)
	g.txBoundary.UpdateLimitPerAddr(ctx, sdk.AccAddress(msg.DelegatorAddress))
	if boundary.TxLimit == 0 {
		return nil
	} else if g.txBoundary.GetLimitPerAddr(ctx, sdk.AccAddress(msg.DelegatorAddress)).DelegateCount > boundary.TxLimit {
		return fmt.Errorf("delegate tx denied, excess tx limit")
	}
	g.txBoundary.IncrementDelegateCount(ctx, sdk.AccAddress(msg.DelegatorAddress))

	return nil
}

func (g StakingPermissionDecorator) validRedelegateMsg(ctx sdk.Context, msg *stakingtypes.MsgBeginRedelegate) error {
	boundary := g.txBoundary.GetRedelegateBoundary(ctx)
	g.txBoundary.UpdateLimitPerAddr(ctx, sdk.AccAddress(msg.DelegatorAddress))
	if boundary.TxLimit == 0 {
		return nil
	} else if g.txBoundary.GetLimitPerAddr(ctx, sdk.AccAddress(msg.DelegatorAddress)).ReledegateCount > boundary.TxLimit {
		return fmt.Errorf("redelegate tx denied, excess tx limit")
	}
	g.txBoundary.IncrementRedelegateCount(ctx, sdk.AccAddress(msg.DelegatorAddress))
	return nil
}

func (g StakingPermissionDecorator) validAuthz(ctx sdk.Context, execMsg *authz.MsgExec) error {
	for _, v := range execMsg.Msgs {
		var innerMsg sdk.Msg
		if err := g.cdc.UnpackAny(v, &innerMsg); err != nil {
			return errorsmod.Wrap(err, "cannot unmarshal authz exec msgs")
		}
		if err := g.ValidateStakingMsg(ctx, innerMsg); err != nil {
			return err
		}
	}

	return nil
}
