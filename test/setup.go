package test

import (
	"testing"

	tmdb "github.com/cometbft/cometbft-db"
	"github.com/cometbft/cometbft/libs/log"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	porttypes "github.com/cosmos/ibc-go/v7/modules/core/05-port/types"
	"github.com/golang/mock/gomock"
	"github.com/notional-labs/banksy/v2/test/mock"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/keeper"
	"github.com/notional-labs/banksy/v2/x/transfermiddleware/types"
	"github.com/stretchr/testify/require"
)

func NewTestSetup(t *testing.T, ctl *gomock.Controller) *Setup {
	t.Helper()
	initializer := newInitializer()
	bankKeeperMock := mock.NewMockBankKeeper(ctl)
	transferKeeperMock := mock.NewMockTransferKeeper(ctl)
	distributionKeeperMock := mock.NewMockDistributionKeeper(ctl)
	ibcModuleMock := mock.NewMockIBCModule(ctl)
	ics4WrapperMock := mock.NewMockICS4Wrapper(ctl)

	paramsKeeper := initializer.paramsKeeper()
	transfermiddlewareKeeper := initializer.transfermiddlewareKeeper(transferKeeperMock, bankKeeperMock)
	// routerModule := initializer.routerModule(routerKeeper)

	require.NoError(t, initializer.StateStore.LoadLatestVersion())

	return &Setup{
		Initializer: initializer,

		Keepers: &testKeepers{
			ParamsKeeper: &paramsKeeper,
			RouterKeeper: &transfermiddlewareKeeper,
		},

		Mocks: &testMocks{
			TransferKeeperMock:     transferKeeperMock,
			BankKeeperMock:         bankKeeperMock,
			DistributionKeeperMock: distributionKeeperMock,
			IBCModuleMock:          ibcModuleMock,
			ICS4WrapperMock:        ics4WrapperMock,
		},

		IBCMiddleware: initializer.IBCMiddleware(ibcModuleMock, transfermiddlewareKeeper),
	}
}

type Setup struct {
	Initializer initializer

	Keepers *testKeepers
	Mocks   *testMocks

	IBCMiddleware transfermiddleware.IBCMiddleware
}

type testKeepers struct {
	ParamsKeeper *paramskeeper.Keeper
	RouterKeeper *keeper.Keeper
}

type testMocks struct {
	TransferKeeperMock     *mock.MockTransferKeeper
	BankKeeperMock         *mock.MockBankKeeper
	DistributionKeeperMock *mock.MockDistributionKeeper
	IBCModuleMock          *mock.MockIBCModule
	ICS4WrapperMock        *mock.MockICS4Wrapper
}

type initializer struct {
	DB         *tmdb.MemDB
	StateStore store.CommitMultiStore
	Ctx        sdk.Context
	Marshaler  codec.Codec
	Amino      *codec.LegacyAmino
}

// Create an initializer with in memory database and default codecs
func newInitializer() initializer {
	logger := log.TestingLogger()
	logger.Debug("initializing test setup")

	db := tmdb.NewMemDB()
	stateStore := store.NewCommitMultiStore(db)

	ctx := sdk.NewContext(stateStore, tmproto.Header{}, false, logger)
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	amino := codec.NewLegacyAmino()
	marshaler := codec.NewProtoCodec(interfaceRegistry)

	return initializer{
		DB:         db,
		StateStore: stateStore,
		Ctx:        ctx,
		Marshaler:  marshaler,
		Amino:      amino,
	}
}

func (i initializer) paramsKeeper() paramskeeper.Keeper {
	storeKey := sdk.NewKVStoreKey(paramstypes.StoreKey)
	transientStoreKey := sdk.NewTransientStoreKey(paramstypes.TStoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)
	i.StateStore.MountStoreWithDB(transientStoreKey, storetypes.StoreTypeTransient, i.DB)

	paramsKeeper := paramskeeper.NewKeeper(i.Marshaler, i.Amino, storeKey, transientStoreKey)

	return paramsKeeper
}

func (i initializer) transfermiddlewareKeeper(
	transferKeeper types.TransferKeeper,
	bankKeeper types.BankKeeper,
) keeper.Keeper {
	storeKey := sdk.NewKVStoreKey(types.StoreKey)
	i.StateStore.MountStoreWithDB(storeKey, storetypes.StoreTypeIAVL, i.DB)

	transfermiddlewareKeeper := keeper.NewKeeper(
		storeKey,
		moduletestutil.MakeTestEncodingConfig().Codec,
		transferKeeper,
		bankKeeper,
	)

	return transfermiddlewareKeeper
}

func (i initializer) IBCMiddleware(app porttypes.IBCModule, k keeper.Keeper) transfermiddleware.IBCMiddleware {
	return transfermiddleware.NewIBCMiddleware(app, k)
}
