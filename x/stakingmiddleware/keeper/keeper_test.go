package keeper_test

import (
	"testing"
	"time"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	moduletestutil "github.com/cosmos/cosmos-sdk/types/module/testutil"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/notional-labs/composable/v6/app"
	"github.com/notional-labs/composable/v6/app/helpers"
	stakingmiddlewarekeeper "github.com/notional-labs/composable/v6/x/stakingmiddleware/keeper"
	stakingmiddlewaretypes "github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
	"github.com/stretchr/testify/suite"
)

type KeeperTestSuite struct {
	suite.Suite

	ctx sdk.Context
	// querier sdk.Querier
	app       *app.ComposableApp
	msgServer stakingmiddlewaretypes.MsgServer
}

func (suite *KeeperTestSuite) SetupTest() {
	suite.app = helpers.SetupComposableAppWithValSet(suite.T())
	suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: 1, ChainID: "centauri-1", Time: time.Now().UTC()})
	encCfg := moduletestutil.MakeTestEncodingConfig()
	key := suite.app.GetKey(stakingmiddlewaretypes.StoreKey)
	keeper := stakingmiddlewarekeeper.NewKeeper(
		encCfg.Codec,
		key,
		suite.app.AccountKeeper,
		suite.app.BankKeeper,
		authtypes.NewModuleAddress(govtypes.ModuleName).String(),
	)
	keeper.RegisterKeepers(suite.app.StakingKeeper)
	suite.msgServer = stakingmiddlewarekeeper.NewMsgServerImpl(keeper)
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

var (
	newParams = stakingmiddlewaretypes.Params{
		BlocksPerEpoch:                           10,
		AllowUnbondAfterEpochProgressBlockNumber: 7,
	}
	failParams = stakingmiddlewaretypes.Params{
		BlocksPerEpoch:                           3,
		AllowUnbondAfterEpochProgressBlockNumber: 3,
	}
	failAllowParams = stakingmiddlewaretypes.Params{
		BlocksPerEpoch:                           6,
		AllowUnbondAfterEpochProgressBlockNumber: 10,
	}
)

func (suite *KeeperTestSuite) TestSetParams() {
	for _, tc := range []struct {
		desc           string
		expectedParams stakingmiddlewaretypes.Params
		malleate       func() error
		shouldErr      bool
		expectedErr    string
	}{
		{
			desc:           "Case success",
			expectedParams: newParams,
			malleate: func() error {
				return suite.app.StakingMiddlewareKeeper.SetParams(suite.ctx, newParams)
			},
			shouldErr: false,
		},
		{
			desc:           "Case fail: BlocksPerEpoch < 5",
			expectedParams: failParams,
			malleate: func() error {
				return suite.app.StakingMiddlewareKeeper.SetParams(suite.ctx, failParams)
			},
			shouldErr:   true,
			expectedErr: "BlocksPerEpoch must be greater than or equal to 5",
		},
		{
			desc:           "Case fail: BlocksPerEpoch < AllowUnbondAfterEpochProgressBlockNumber",
			expectedParams: failAllowParams,
			malleate: func() error {
				return suite.app.StakingMiddlewareKeeper.SetParams(suite.ctx, failAllowParams)
			},
			shouldErr:   true,
			expectedErr: "AllowUnbondAfterEpochProgressBlockNumber must be less than or equal to BlocksPerEpoch",
		},
	} {
		tc := tc
		suite.Run(tc.desc, func() {
			suite.SetupTest()
			err := tc.malleate()
			if !tc.shouldErr {
				res := suite.app.StakingMiddlewareKeeper.GetParams(suite.ctx)
				suite.Equal(res, tc.expectedParams)
			} else {
				suite.Equal(err.Error(), tc.expectedErr)
			}
		})
	}
}
