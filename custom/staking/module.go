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

// func (am AppModule) BeginBlock(ctx sdk.Context, _abc abcitype.RequestBeginBlock) {
// 	//Define the logic around the batching.
// 	am.AppModule.BeginBlock(ctx, _abc)
// }

func (am AppModule) EndBlock(ctx sdk.Context, _abc abcitype.RequestEndBlock) []abcitype.ValidatorUpdate {
	//Define the logic around the batching.
	//TODO!!!

	println("EndBlock Custom Staking Module")

	delegations := am.keeper.Stakingmiddleware.DequeueAllDelegation(ctx)
	println("Delegations: ", delegations)
	println("Delegations len: ", len(delegations))
	//for delegations print the delegator address and the validator address
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

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			"DequeueAllDelegation",
			// sdk.NewAttribute(sdk.AttributeKeyAmount, balances.String()),
			// sdk.NewAttribute(types.AttributeKeyValidator, dvPair.ValidatorAddress),
			// sdk.NewAttribute(types.AttributeKeyDelegator, dvPair.DelegatorAddress),
		),
	)

	return am.AppModule.EndBlock(ctx, _abc)
}
