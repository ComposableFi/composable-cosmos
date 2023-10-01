package keeper_test

import (
	"encoding/binary"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"

	"github.com/notional-labs/centauri/v5/app"
	"github.com/notional-labs/centauri/v5/app/helpers"
	"github.com/notional-labs/centauri/v5/x/ratelimit/keeper"
	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

var (
	sampleRateLimitA = types.MsgAddRateLimit{
		Denom:              "denomA",
		ChannelID:          "channel-0",
		MaxPercentSend:     sdkmath.NewInt(10),
		MaxPercentRecv:     sdkmath.NewInt(10),
		MinRateLimitAmount: sdkmath.NewInt(1_000_000),
		DurationHours:      uint64(1),
	}
	sampleRateLimitB = types.MsgAddRateLimit{
		Denom:              "denomB",
		ChannelID:          "channel-0",
		MaxPercentSend:     sdkmath.NewInt(20),
		MaxPercentRecv:     sdkmath.NewInt(20),
		MinRateLimitAmount: sdkmath.NewInt(1_000_000),
		DurationHours:      uint64(1),
	}
	sampleRateLimitC = types.MsgAddRateLimit{
		Denom:              "denomB",
		ChannelID:          "channel-1",
		MaxPercentSend:     sdkmath.NewInt(50),
		MaxPercentRecv:     sdkmath.NewInt(50),
		MinRateLimitAmount: sdkmath.NewInt(5_000_000),
		DurationHours:      uint64(5),
	}
	sampleRateLimitD = types.MsgAddRateLimit{
		Denom:              "denomC",
		ChannelID:          "channel-2",
		MaxPercentSend:     sdkmath.NewInt(80),
		MaxPercentRecv:     sdkmath.NewInt(80),
		MinRateLimitAmount: sdkmath.NewInt(10_000_000),
		DurationHours:      uint64(10),
	}
)

type KeeperTestSuite struct {
	suite.Suite

	app       *app.CentauriApp
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   types.QueryServer
	msgServer types.MsgServer

	coordinator *ibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *ibctesting.TestChain
	chainB *ibctesting.TestChain
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(KeeperTestSuite))
}

func (s *KeeperTestSuite) SetupTest() {
	s.app = helpers.SetupCentauriAppWithValSet(s.T())
	s.ctx = s.app.BaseApp.NewContext(false, tmproto.Header{
		Height:  1,
		ChainID: "centauri-1",
		Time:    time.Now().UTC(),
	})
	s.keeper = s.app.RatelimitKeeper
	s.querier = keeper.NewQueryServer(s.keeper)
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)

	// Creates a coordinator with 2 test chains
	s.coordinator = ibctesting.NewCoordinator(s.T(), 2)
	s.chainA = s.coordinator.GetChain(ibctesting.GetChainID(1))
	s.chainB = s.coordinator.GetChain(ibctesting.GetChainID(2))

	// Commit some blocks so that QueryProof returns valid proof (cannot return valid query if height <= 1)
	s.coordinator.CommitNBlocks(s.chainA, 2)
	s.coordinator.CommitNBlocks(s.chainB, 2)
}

func (s *KeeperTestSuite) SetupSampleRateLimits(rateLimits ...types.MsgAddRateLimit) {
	for _, rateLimit := range rateLimits {
		s.addRateLimit(
			rateLimit.Denom,
			rateLimit.ChannelID,
			rateLimit.MaxPercentSend,
			rateLimit.MaxPercentRecv,
			rateLimit.MinRateLimitAmount,
			rateLimit.DurationHours,
		)
	}
}

//
// Below are helper functions to write test code easily
//

func (s *KeeperTestSuite) addr(addrNum int) sdk.AccAddress {
	addr := make(sdk.AccAddress, 20)
	binary.PutVarint(addr, int64(addrNum))
	return addr
}

func (s *KeeperTestSuite) fundAddr(addr sdk.AccAddress, amt sdk.Coins) {
	s.T().Helper()
	err := s.app.BankKeeper.MintCoins(s.ctx, minttypes.ModuleName, amt)
	s.Require().NoError(err)
	err = s.app.BankKeeper.SendCoinsFromModuleToAccount(s.ctx, minttypes.ModuleName, addr, amt)
	s.Require().NoError(err)
}

// addRateLimit is a convenient method to add new RateLimit without the need of authority.
func (s *KeeperTestSuite) addRateLimit(
	denom string,
	channelID string,
	maxPercentSend sdkmath.Int,
	maxPercentRecv sdkmath.Int,
	minRateLimitAmount sdkmath.Int,
	durationHours uint64,
) {
	s.T().Helper()

	// Add new RateLimit requires total supply of the given denom
	s.fundAddr(s.addr(0), sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(100_000_000_000))))

	err := s.keeper.AddRateLimit(s.ctx, &types.MsgAddRateLimit{
		Authority:          "",
		Denom:              denom,
		ChannelID:          channelID,
		MaxPercentSend:     maxPercentSend,
		MaxPercentRecv:     maxPercentRecv,
		MinRateLimitAmount: minRateLimitAmount,
		DurationHours:      durationHours,
	})
	s.Require().NoError(err)
}
