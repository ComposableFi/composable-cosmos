package ratelimit_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	customibctesting "github.com/notional-labs/centauri/v3/app/ibctesting"
	ratelimittypes "github.com/notional-labs/centauri/v3/x/ratelimit/types"
	"github.com/stretchr/testify/suite"
)

type RateLimitTestSuite struct {
	suite.Suite

	coordinator *customibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *customibctesting.TestChain
	chainB *customibctesting.TestChain
	chainC *customibctesting.TestChain
}

func (suite *RateLimitTestSuite) SetupTest() {
	suite.coordinator = customibctesting.NewCoordinator(suite.T(), 4)
	suite.chainA = suite.coordinator.GetChain(customibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(customibctesting.GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(customibctesting.GetChainID(3))
}

func NewTransferPath(chainA, chainB *customibctesting.TestChain) *customibctesting.Path {
	path := customibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = customibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = customibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = transfertypes.Version
	path.EndpointB.ChannelConfig.Version = transfertypes.Version

	return path
}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimitTestSuite))
}

func (suite *RateLimitTestSuite) TestReceiveIBCToken() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		// when transfer via sdk transfer from A (module) -> B (contract)
		ibcDenom                   = "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878"
		nativeDenom                = "ppica"
		nativeTokenSendOnChainA    = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
		nativeTokenReceiveOnChainB = sdk.NewCoin(nativeDenom, transferAmount)
		timeoutHeight              = clienttypes.NewHeight(1, 110)
		expChainABalanceDiff       = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
	)

	suite.SetupTest() // reset

	path := NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	// Add parachain token info
	chainBtransMiddlewareKeeper := suite.chainB.TransferMiddleware()
	err := chainBtransMiddlewareKeeper.AddParachainIBCInfo(suite.chainB.GetContext(), ibcDenom, path.EndpointB.ChannelID, nativeDenom, sdk.DefaultBondDenom)
	suite.Require().NoError(err)

	originalChainABalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	originalChainBBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())

	msg := transfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nativeTokenSendOnChainA, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
	_, err = suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err)
	suite.Require().NoError(err, path.EndpointB.UpdateClient())

	// then
	suite.Require().Equal(1, len(suite.chainA.PendingSendPackets))
	suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

	// and when relay to chain B and handle Ack on chain A
	err = suite.coordinator.RelayAndAckPendingPackets(path)
	suite.Require().NoError(err)

	// then
	suite.Require().Equal(0, len(suite.chainA.PendingSendPackets))
	suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

	// and source chain balance was decreased
	newChainABalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	suite.Require().Equal(originalChainABalance.Sub(expChainABalanceDiff), newChainABalance)

	// and dest chain balance contains voucher
	expBalance := originalChainBBalance.Add(nativeTokenReceiveOnChainB)
	gotBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())
	suite.Require().Equal(expBalance, gotBalance)

	// add rate limit
	chainBRateLimitKeeper := suite.chainB.RateLimit()
	err = chainBRateLimitKeeper.AddRateLimit(suite.chainB.GetContext(), &ratelimittypes.MsgAddRateLimit{})
}
