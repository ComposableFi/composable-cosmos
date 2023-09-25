package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"

	"github.com/notional-labs/centauri/v5/app"
	"github.com/notional-labs/centauri/v5/app/helpers"
	"github.com/notional-labs/centauri/v5/x/ratelimit/keeper"
	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

type KeeperTestSuite struct {
	suite.Suite

	app       *app.CentauriApp
	ctx       sdk.Context
	keeper    keeper.Keeper
	querier   types.QueryServer
	msgServer types.MsgServer
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
}

// addRateLimit is a convenient method to add new RateLimit without the need of authority.
func (s *KeeperTestSuite) addRateLimit(
	denom string,
	channelID string,
	maxPercentSend sdkmath.Int,
	maxPercentRecv sdkmath.Int,
	MinRateLimitAmount sdkmath.Int,
	DurationHours uint64,
) {
	s.T().Helper()

	err := s.keeper.AddRateLimit(s.ctx, &types.MsgAddRateLimit{
		Authority:          "",
		Denom:              denom,
		ChannelID:          channelID,
		MaxPercentSend:     maxPercentSend,
		MaxPercentRecv:     maxPercentRecv,
		MinRateLimitAmount: MinRateLimitAmount,
		DurationHours:      DurationHours,
	})
	s.Require().NoError(err)
}
