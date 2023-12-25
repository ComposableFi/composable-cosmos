package ante

import (
	tfmwkeeper "github.com/notional-labs/composable/v6/x/transfermiddleware/keeper"
	txboundaryante "github.com/notional-labs/composable/v6/x/txboundary/ante"
	txboundarykeeper "github.com/notional-labs/composable/v6/x/txboundary/keeper"

	"github.com/cosmos/cosmos-sdk/codec"
	servertypes "github.com/cosmos/cosmos-sdk/server/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/ante"
	"github.com/cosmos/cosmos-sdk/x/auth/signing"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	ibcante "github.com/cosmos/ibc-go/v7/modules/core/ante"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
)

// Link to default ante handler used by cosmos sdk:
// https://github.com/cosmos/cosmos-sdk/blob/v0.43.0/x/auth/ante/ante.go#L41
func NewAnteHandler(
	options servertypes.AppOptions,
	ak ante.AccountKeeper,
	bk authtypes.BankKeeper,
	feegrantKeeper ante.FeegrantKeeper,
	txFeeChecker ante.TxFeeChecker,
	sigGasConsumer ante.SignatureVerificationGasConsumer,
	signModeHandler signing.SignModeHandler,
	channelKeeper *ibckeeper.Keeper,
	tfmwKeeper tfmwkeeper.Keeper,
	txBoundaryKeeper txboundarykeeper.Keeper,
	codec codec.BinaryCodec,
) sdk.AnteHandler {
	return sdk.ChainAnteDecorators(
		ante.NewSetUpContextDecorator(), //  // outermost AnteDecorator. SetUpContext must be called first
		ante.NewValidateBasicDecorator(),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		ante.NewDeductFeeDecorator(ak, bk, feegrantKeeper, txFeeChecker),
		ante.NewTxTimeoutHeightDecorator(),
		ante.NewValidateMemoDecorator(ak),
		ante.NewConsumeGasForTxSizeDecorator(ak),
		NewIBCPermissionDecorator(codec, tfmwKeeper),
		txboundaryante.NewStakingPermissionDecorator(codec, txBoundaryKeeper),
		ante.NewSetPubKeyDecorator(ak), // SetPubKeyDecorator must be called before all signature verification decorators
		ante.NewValidateSigCountDecorator(ak),
		ante.NewSigGasConsumeDecorator(ak, sigGasConsumer),
		ante.NewSigVerificationDecorator(ak, signModeHandler),
		ante.NewIncrementSequenceDecorator(ak),
		ibcante.NewRedundantRelayDecorator(channelKeeper),
	)
}
