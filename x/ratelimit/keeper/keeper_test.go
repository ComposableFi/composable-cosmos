package keeper_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

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
