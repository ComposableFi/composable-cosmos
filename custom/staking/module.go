package bank

import (
	"fmt"

	abcitype "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingmodule "github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/staking/exported"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	// custombankkeeper "github.com/notional-labs/composable/v6/custom/bank/keeper"
	sdk "github.com/cosmos/cosmos-sdk/types"
	customstakingkeeper "github.com/notional-labs/composable/v6/custom/staking/keeper"
)

// AppModule wraps around the bank module and the bank keeper to return the right total supply
type AppModule struct {
	stakingmodule.AppModule
	keeper    customstakingkeeper.Keeper
	subspace  exported.Subspace
	msgServer stakingtypes.MsgServer
}

// NewAppModule creates a new AppModule object
func NewAppModule(cdc codec.Codec, keeper customstakingkeeper.Keeper, accountKeeper stakingtypes.AccountKeeper, bankKeeper stakingtypes.BankKeeper, ss exported.Subspace) AppModule {
	stakingModule := stakingmodule.NewAppModule(cdc, &keeper.Keeper, accountKeeper, bankKeeper, ss)
	return AppModule{
		AppModule: stakingModule,
		keeper:    keeper,
		subspace:  ss,
		msgServer: stakingkeeper.NewMsgServerImpl(&keeper.Keeper),
	}
}

// RegisterServices registers module services.
// NOTE: Overriding this method as not doing so will cause a panic
// when trying to force this custom keeper into a bankkeeper.BaseKeeper
func (am AppModule) RegisterServices(cfg module.Configurator) {
	// types.RegisterMsgServer(cfg.MsgServer(), keeper.NewMsgServerImpl(&am.keeper))
	stakingtypes.RegisterMsgServer(cfg.MsgServer(), customstakingkeeper.NewMsgServerImpl(am.keeper.Keeper, am.keeper))
	querier := stakingkeeper.Querier{Keeper: &am.keeper.Keeper}
	stakingtypes.RegisterQueryServer(cfg.QueryServer(), querier)

	m := stakingkeeper.NewMigrator(&am.keeper.Keeper, am.subspace)
	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 1, m.Migrate1to2); err != nil {
		panic(fmt.Sprintf("failed to migrate x/staking from version 1 to 2: %v", err))
	}

	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 2, m.Migrate2to3); err != nil {
		panic(fmt.Sprintf("failed to migrate x/staking from version 2 to 3: %v", err))
	}

	if err := cfg.RegisterMigration(stakingtypes.ModuleName, 3, m.Migrate3to4); err != nil {
		panic(fmt.Sprintf("failed to migrate x/staking from version 3 to 4: %v", err))
	}
}

func (am AppModule) EndBlock(ctx sdk.Context, _abc abcitype.RequestEndBlock) []abcitype.ValidatorUpdate {

	return EndBlocker(ctx, &am.keeper)

	println("EndBlock Custom Staking Module")
	params := am.keeper.Stakingmiddleware.GetParams(ctx)
	println("BlocksPerEpoch: ", params.BlocksPerEpoch)
	println("Height: ", _abc.Height)

	should_execute_batch := (_abc.Height % int64(params.BlocksPerEpoch)) == 0
	if should_execute_batch {
		println("Should batch delegation to be executed at block: ", _abc.Height)

		delegations := am.keeper.Stakingmiddleware.DequeueAllDelegation(ctx)
		println("Delegations: ", delegations)
		println("Delegations len: ", len(delegations))
		for _, delegation := range delegations {
			println("Delegator Address: ", delegation.DelegatorAddress)
			println("Validator Address: ", delegation.ValidatorAddress)
			fmt.Println("Amount", delegation.Amount.Amount)

			msgDelegate := stakingtypes.MsgDelegate{DelegatorAddress: delegation.DelegatorAddress, ValidatorAddress: delegation.ValidatorAddress, Amount: delegation.Amount}
			_, err := am.msgServer.Delegate(ctx, &msgDelegate)
			if err != nil {
				println("Error for Delegator Address: ", delegation.DelegatorAddress)
			}
		}

		beginredelegations := am.keeper.Stakingmiddleware.DequeueAllRedelegation(ctx)
		println("BeginRedelegations: ", beginredelegations)
		println("BeginRedelegations len: ", len(beginredelegations))
		for _, redelegation := range beginredelegations {
			println("Delegator Address: ", redelegation.DelegatorAddress)
			println("Validator Address: ", redelegation.ValidatorSrcAddress)

			msg_redelegation := stakingtypes.MsgBeginRedelegate{DelegatorAddress: redelegation.DelegatorAddress, ValidatorSrcAddress: redelegation.ValidatorSrcAddress, ValidatorDstAddress: redelegation.ValidatorDstAddress, Amount: redelegation.Amount}
			_, err := am.msgServer.BeginRedelegate(ctx, &msg_redelegation)
			if err != nil {
				println("Error for Delegator Address: ", msg_redelegation.DelegatorAddress)
			}
		}

		undelegations := am.keeper.Stakingmiddleware.DequeueAllUndelegation(ctx)
		println("Undelegation: ", beginredelegations)
		println("Undelegation len: ", len(beginredelegations))
		for _, undelegation := range undelegations {
			println("Undelegation Delegator Address: ", undelegation.DelegatorAddress)
			println("Undelegation Validator Address: ", undelegation.ValidatorAddress)

			msg_undelegate := stakingtypes.MsgUndelegate{DelegatorAddress: undelegation.DelegatorAddress, ValidatorAddress: undelegation.ValidatorAddress, Amount: undelegation.Amount}
			_, err := am.msgServer.Undelegate(ctx, &msg_undelegate)
			if err != nil {
				println("Error for Delegator Address: ", msg_undelegate.DelegatorAddress)
			}
		}

		cancel_unbonding_delegations := am.keeper.Stakingmiddleware.DequeueAllCancelUnbondingDelegation(ctx)
		println("Cancel Unbonding Delegations: ", cancel_unbonding_delegations)
		println("Cancel Ubonding Delegations len: ", len(cancel_unbonding_delegations))
		for _, cancel_unbonding_delegation := range cancel_unbonding_delegations {
			println("Cancel Unbonding Delegation  Delegator Address: ", cancel_unbonding_delegation.DelegatorAddress)
			println("Cancel Unbonding Delegations Validator Address: ", cancel_unbonding_delegation.ValidatorAddress)

			msg_cancle_unbonding_delegation := stakingtypes.MsgCancelUnbondingDelegation{DelegatorAddress: cancel_unbonding_delegation.DelegatorAddress, ValidatorAddress: cancel_unbonding_delegation.ValidatorAddress, Amount: cancel_unbonding_delegation.Amount, CreationHeight: cancel_unbonding_delegation.CreationHeight}
			_, err := am.msgServer.CancelUnbondingDelegation(ctx, &msg_cancle_unbonding_delegation)
			if err != nil {
				println("Error for Delegator Address: ", msg_cancle_unbonding_delegation.DelegatorAddress)
			}
		}
	}

	return am.AppModule.EndBlock(ctx, _abc)
}
