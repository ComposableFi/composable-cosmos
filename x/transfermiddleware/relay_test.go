package transfermiddleware_test

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/cometbft/cometbft/crypto"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	customibctesting "github.com/notional-labs/banksy/v2/app/ibctesting"
	"github.com/stretchr/testify/suite"
)

type PacketMetadata struct {
	Forward *ForwardMetadata `json:"forward"`
}

type ForwardMetadata struct {
	Receiver string        `json:"receiver,omitempty"`
	Port     string        `json:"port,omitempty"`
	Channel  string        `json:"channel,omitempty"`
	Timeout  time.Duration `json:"timeout,omitempty"`
	Retries  *uint8        `json:"retries,omitempty"`
	Next     string        `json:"next,omitempty"`
}

// TODO: use testsuite here.
type TransferMiddlewareTestSuite struct {
	suite.Suite

	coordinator *customibctesting.Coordinator

	// testing chains used for convenience and readability
	chainA *customibctesting.TestChain
	chainB *customibctesting.TestChain
	chainC *customibctesting.TestChain
}

func (suite *TransferMiddlewareTestSuite) SetupTest() {
	suite.coordinator = customibctesting.NewCoordinator(suite.T(), 4)
	suite.chainA = suite.coordinator.GetChain(customibctesting.GetChainID(1))
	suite.chainB = suite.coordinator.GetChain(customibctesting.GetChainID(2))
	suite.chainC = suite.coordinator.GetChain(customibctesting.GetChainID(3))

}

func NewTransferPath(chainA, chainB *customibctesting.TestChain) *customibctesting.Path {
	path := customibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig.PortID = customibctesting.TransferPort
	path.EndpointB.ChannelConfig.PortID = customibctesting.TransferPort
	path.EndpointA.ChannelConfig.Version = ibctransfertypes.Version
	path.EndpointB.ChannelConfig.Version = ibctransfertypes.Version

	return path

}

func TestKeeperTestSuite(t *testing.T) {
	suite.Run(t, new(TransferMiddlewareTestSuite))
}

var keyCounter uint64

// we need to make this deterministic (same every test run), as encoded address size and thus gas cost,
// depends on the actual bytes (due to ugly CanonicalAddress encoding)
func keyPubAddr() (crypto.PrivKey, crypto.PubKey, sdk.AccAddress) {
	keyCounter++
	seed := make([]byte, 8)
	binary.BigEndian.PutUint64(seed, keyCounter)

	key := ed25519.GenPrivKeyFromSecret(seed)
	pub := key.PubKey()
	addr := sdk.AccAddress(pub.Address())
	return key, pub, addr
}

func RandomAccountAddress(_ testing.TB) sdk.AccAddress {
	_, _, addr := keyPubAddr()
	return addr
}

func RandomBech32AccountAddress(t testing.TB) string {
	return RandomAccountAddress(t).String()
}

func (suite *TransferMiddlewareTestSuite) TestTransferWithPFM_ErrorAck() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		timeoutHeight  = clienttypes.NewHeight(1, 110)
		pathAtoB       *customibctesting.Path
		pathBtoC       *customibctesting.Path
		ibcDenom       = "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878"
	)

	suite.SetupTest()
	pathAtoB = NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(pathAtoB)
	pathBtoC = NewTransferPath(suite.chainB, suite.chainC)
	suite.coordinator.Setup(pathBtoC)
	// Add parachain token info
	chainBtransMiddleware := suite.chainB.TransferMiddleware()
	err := chainBtransMiddleware.AddParachainIBCInfo(suite.chainB.GetContext(), ibcDenom, pathAtoB.EndpointB.ChannelID, sdk.DefaultBondDenom)
	suite.Require().NoError(err)

	params := transfertypes.Params{
		SendEnabled:    true,
		ReceiveEnabled: false,
	}

	// set send params
	suite.chainC.GetTestSupport().TransferKeeper().SetParams(suite.chainC.GetContext(), params)
	senderAOriginalBalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())

	testAcc := RandomAccountAddress(suite.T())
	timeOut := 10 * time.Minute
	retries := uint8(0)
	// Build MEMOtransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	memo := PacketMetadata{
		Forward: &ForwardMetadata{
			Receiver: testAcc.String(),
			Port:     pathBtoC.EndpointA.ChannelConfig.PortID,
			Channel:  pathBtoC.EndpointA.ChannelID,
			Timeout:  timeOut,
			Retries:  &retries,
			Next:     "",
		},
	}
	memo_marshalled, err := json.Marshal(&memo)
	suite.Require().NoError(err)

	msg := ibctransfertypes.NewMsgTransfer(
		pathAtoB.EndpointA.ChannelConfig.PortID,
		pathAtoB.EndpointA.ChannelID,
		sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
		suite.chainA.SenderAccount.GetAddress().String(),
		testAcc.String(),
		timeoutHeight,
		0,
		string(memo_marshalled),
	)
	_, err = suite.chainA.SendMsgs(msg)
	suite.Require().NoError(err)
	suite.Require().NoError(err, pathAtoB.EndpointB.UpdateClient())

	// then
	suite.Require().Equal(1, len(suite.chainA.PendingSendPackets))
	suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))
	// relay packet
	sendingPacket := suite.chainA.PendingSendPackets[0]
	suite.coordinator.IncrementTime()
	suite.coordinator.CommitBlock(suite.chainA)
	err = pathAtoB.EndpointB.UpdateClient()
	suite.Require().NoError(err)

	err = pathAtoB.EndpointB.RecvPacket(sendingPacket)
	suite.Require().NoError(err)
	suite.chainA.PendingSendPackets = nil

	// Check after first Hop
	senderABalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	suite.Require().Equal(senderAOriginalBalance.Sub(sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)), senderABalance)

	escrowIbcDenomAddress := transfertypes.GetEscrowAddress(pathAtoB.EndpointB.ChannelConfig.PortID, pathAtoB.EndpointB.ChannelID)
	escrowIbcDenomAddressBalance := suite.chainB.AllBalances(escrowIbcDenomAddress)
	suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(ibcDenom, transferAmount)), escrowIbcDenomAddressBalance)

	escrowNativeDenomAddress := transfertypes.GetEscrowAddress(pathBtoC.EndpointA.ChannelConfig.PortID, pathBtoC.EndpointA.ChannelID)
	escrowNativeDenomAddressBalance := suite.chainB.AllBalances(escrowNativeDenomAddress)
	suite.Require().Equal(sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)), escrowNativeDenomAddressBalance)

	// then should have a packet from B to C
	suite.Require().Equal(1, len(suite.chainB.PendingSendPackets))
	suite.Require().Equal(0, len(suite.chainC.PendingSendPackets))

	// relay packet
	sendingPacket = suite.chainB.PendingSendPackets[0]
	suite.coordinator.IncrementTime()
	suite.coordinator.CommitBlock(suite.chainB)
	err = pathBtoC.EndpointB.UpdateClient()
	suite.Require().NoError(err)

	err = pathBtoC.EndpointB.RecvPacket(sendingPacket)
	suite.Require().NoError(err)
	suite.chainB.PendingSendPackets = nil

	// relay ack C to B
	suite.Require().Equal(1, len(suite.chainC.PendingAckPackets))
	ack := suite.chainC.PendingAckPackets[0]
	suite.coordinator.IncrementTime()
	suite.coordinator.CommitBlock(suite.chainC)
	err = pathBtoC.EndpointA.UpdateClient()
	suite.Require().NoError(err)
	// relay failed ack
	err = pathBtoC.EndpointA.AcknowledgePacket(ack.Packet, ack.Ack)
	suite.Require().NoError(err)
	suite.chainC.PendingAckPackets = nil

	// relay ack B to A
	suite.Require().Equal(1, len(suite.chainB.PendingAckPackets))
	ack = suite.chainB.PendingAckPackets[0]
	suite.coordinator.IncrementTime()
	suite.coordinator.CommitBlock(suite.chainB)
	err = pathAtoB.EndpointA.UpdateClient()
	suite.Require().NoError(err)

	err = pathAtoB.EndpointA.AcknowledgePacket(ack.Packet, ack.Ack)
	suite.Require().NoError(err)
	suite.chainB.PendingAckPackets = nil

	escrowIbcDenomAddress = transfertypes.GetEscrowAddress(pathAtoB.EndpointB.ChannelConfig.PortID, pathAtoB.EndpointB.ChannelID)
	escrowIbcDenomAddressBalance = suite.chainB.AllBalances(escrowIbcDenomAddress)
	suite.Require().Empty(escrowIbcDenomAddressBalance)

	escrowNativeDenomAddress = transfertypes.GetEscrowAddress(pathBtoC.EndpointA.ChannelConfig.PortID, pathBtoC.EndpointA.ChannelID)
	escrowNativeDenomAddressBalance = suite.chainB.AllBalances(escrowNativeDenomAddress)
	suite.Require().Empty(escrowNativeDenomAddressBalance)

	balance := suite.chainB.AllBalances(testAcc)
	suite.Require().Empty(balance)

	balance = suite.chainC.AllBalances(testAcc)
	suite.Require().Empty(balance)

	balance = suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	suite.Require().Equal(senderAOriginalBalance, balance)
}

func (suite *TransferMiddlewareTestSuite) TestTransferWithPFM() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		// when transfer via sdk transfer from A (module) -> B (contract)
		timeoutHeight = clienttypes.NewHeight(1, 110)
		pathAtoB      *customibctesting.Path
		pathBtoC      *customibctesting.Path
	)

	testCases := []struct {
		name string
	}{
		{
			"Success case A -> B -> C",
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			pathAtoB = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(pathAtoB)
			pathBtoC = NewTransferPath(suite.chainB, suite.chainC)
			suite.coordinator.Setup(pathBtoC)

			// Add parachain token info
			chainBtransMiddleware := suite.chainB.TransferMiddleware()
			err := chainBtransMiddleware.AddParachainIBCInfo(suite.chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", pathAtoB.EndpointB.ChannelID, sdk.DefaultBondDenom)
			suite.Require().NoError(err)

			testAcc := RandomAccountAddress(suite.T())
			timeOut := 10 * time.Minute
			retries := uint8(0)
			// Build MEMO
			memo := PacketMetadata{
				Forward: &ForwardMetadata{
					Receiver: testAcc.String(),
					Port:     pathBtoC.EndpointA.ChannelConfig.PortID,
					Channel:  pathBtoC.EndpointA.ChannelID,
					Timeout:  timeOut,
					Retries:  &retries,
					Next:     "",
				},
			}
			memo_marshalled, err := json.Marshal(&memo)
			suite.Require().NoError(err)

			intermediaryOriginalBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())

			msg := ibctransfertypes.NewMsgTransfer(
				pathAtoB.EndpointA.ChannelConfig.PortID,
				pathAtoB.EndpointA.ChannelID,
				sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				timeoutHeight,
				0,
				string(memo_marshalled),
			)
			_, err = suite.chainA.SendMsgs(msg)
			suite.Require().NoError(err)
			suite.Require().NoError(err, pathAtoB.EndpointB.UpdateClient())

			// then
			suite.Require().Equal(1, len(suite.chainA.PendingSendPackets))
			suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))
			// relay packet
			sendingPacket := suite.chainA.PendingSendPackets[0]
			suite.coordinator.IncrementTime()
			suite.coordinator.CommitBlock(suite.chainA)
			err = pathAtoB.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			err = pathAtoB.EndpointB.RecvPacket(sendingPacket)
			suite.Require().NoError(err)
			suite.chainA.PendingSendPackets = nil
			// then should have a packet from B to C
			suite.Require().Equal(1, len(suite.chainB.PendingSendPackets))
			suite.Require().Equal(0, len(suite.chainC.PendingSendPackets))

			// relay packet
			sendingPacket = suite.chainB.PendingSendPackets[0]
			suite.coordinator.IncrementTime()
			suite.coordinator.CommitBlock(suite.chainB)
			err = pathBtoC.EndpointB.UpdateClient()
			suite.Require().NoError(err)

			err = pathBtoC.EndpointB.RecvPacket(sendingPacket)
			suite.Require().NoError(err)
			suite.chainB.PendingSendPackets = nil

			// relay ack C to B
			suite.Require().Equal(1, len(suite.chainC.PendingAckPackets))
			ack := suite.chainC.PendingAckPackets[0]
			suite.coordinator.IncrementTime()
			suite.coordinator.CommitBlock(suite.chainC)
			err = pathBtoC.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			err = pathBtoC.EndpointA.AcknowledgePacket(ack.Packet, ack.Ack)
			suite.Require().NoError(err)
			suite.chainC.PendingAckPackets = nil

			// relay ack B to A
			suite.Require().Equal(1, len(suite.chainB.PendingAckPackets))
			ack = suite.chainB.PendingAckPackets[0]
			suite.coordinator.IncrementTime()
			suite.coordinator.CommitBlock(suite.chainB)
			err = pathAtoB.EndpointA.UpdateClient()
			suite.Require().NoError(err)

			err = pathAtoB.EndpointA.AcknowledgePacket(ack.Packet, ack.Ack)
			suite.Require().NoError(err)
			suite.chainB.PendingAckPackets = nil

			intermediaryBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())
			suite.Require().Equal(intermediaryOriginalBalance, intermediaryBalance)
			expDenom := "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878"
			expBalance := sdk.NewCoins(sdk.NewCoin(expDenom, transferAmount))
			balance := suite.chainC.AllBalances(testAcc)
			suite.Require().Equal(expBalance, balance)
		})
	}
}

func (suite *TransferMiddlewareTestSuite) TestSendTransfer() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		// when transfer via sdk transfer from A (module) -> B (contract)
		timeoutHeight = clienttypes.NewHeight(1, 110)
		pathAtoB      *customibctesting.Path
		pathCtoB      *customibctesting.Path
		path          *customibctesting.Path
		srcPort       string
		srcChannel    string
		chain         *customibctesting.TestChain
		expDenom      string
		// pathBtoC      = NewTransferPath(suite.chainB, suite.chainC)
	)

	testCases := []struct {
		name     string
		malleate func()
	}{
		{
			"Receiver is Parachain chain",
			func() {
				path = pathAtoB
				srcPort = pathAtoB.EndpointB.ChannelConfig.PortID
				srcChannel = pathAtoB.EndpointB.ChannelID
				chain = suite.chainA
				expDenom = sdk.DefaultBondDenom
			},
		},
		{
			"Receiver is cosmos chain chain",
			func() {
				path = pathCtoB
				srcPort = pathCtoB.EndpointB.ChannelConfig.PortID
				srcChannel = pathCtoB.EndpointB.ChannelID
				chain = suite.chainC
				expDenom = "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878"
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest()
			pathAtoB = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(pathAtoB)
			pathCtoB = NewTransferPath(suite.chainC, suite.chainB)
			suite.coordinator.Setup(pathCtoB)
			// Add parachain token info
			chainBtransMiddleware := suite.chainB.TransferMiddleware()
			err := chainBtransMiddleware.AddParachainIBCInfo(suite.chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", pathAtoB.EndpointB.ChannelID, sdk.DefaultBondDenom)
			suite.Require().NoError(err)
			// send coin from A to B

			msg := ibctransfertypes.NewMsgTransfer(
				pathAtoB.EndpointA.ChannelConfig.PortID,
				pathAtoB.EndpointA.ChannelID,
				sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
				suite.chainA.SenderAccount.GetAddress().String(),
				suite.chainB.SenderAccount.GetAddress().String(),
				timeoutHeight,
				0,
				"",
			)
			_, err = suite.chainA.SendMsgs(msg)
			suite.Require().NoError(err)
			suite.Require().NoError(err, pathAtoB.EndpointB.UpdateClient())

			// then
			suite.Require().Equal(1, len(suite.chainA.PendingSendPackets))
			suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

			// and when relay to chain A and handle Ack on chain B
			err = suite.coordinator.RelayAndAckPendingPackets(pathAtoB)
			suite.Require().NoError(err)

			// then
			suite.Require().Equal(0, len(suite.chainA.PendingSendPackets))
			suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

			tc.malleate()

			testAcc2 := RandomAccountAddress(suite.T())
			msg = ibctransfertypes.NewMsgTransfer(
				srcPort,
				srcChannel,
				sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(500000)),
				suite.chainB.SenderAccount.GetAddress().String(),
				testAcc2.String(),
				timeoutHeight,
				0,
				"",
			)
			_, err = suite.chainB.SendMsgs(msg)
			suite.Require().NoError(err)
			suite.Require().NoError(err, path.EndpointB.UpdateClient())

			suite.Require().Equal(1, len(suite.chainB.PendingSendPackets))
			suite.Require().Equal(0, len(chain.PendingSendPackets))

			// and when relay to chain B and handle Ack on chain A
			err = suite.coordinator.RelayAndAckPendingPacketsReverse(path)
			suite.Require().NoError(err)

			suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))
			suite.Require().Equal(0, len(chain.PendingSendPackets))

			balance := chain.AllBalances(testAcc2)
			expBalance := sdk.NewCoins(sdk.NewCoin(expDenom, sdk.NewInt(500000)))
			suite.Require().Equal(expBalance, balance)
		})
	}
}

// TODO: use testsuite here.
func (suite *TransferMiddlewareTestSuite) TestOnrecvPacket() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		// when transfer via sdk transfer from A (module) -> B (contract)
		coinToSendToB = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
		timeoutHeight = clienttypes.NewHeight(1, 110)
	)
	var (
		expChainBBalanceDiff sdk.Coin
		path                 = NewTransferPath(suite.chainA, suite.chainB)
	)

	testCases := []struct {
		name                 string
		expChainABalanceDiff sdk.Coin
		malleate             func()
	}{
		{
			"Transfer with no pre-set ParachainIBCTokenInfo",
			sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
			func() {
				expChainBBalanceDiff = ibctransfertypes.GetTransferCoin(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, coinToSendToB.Denom, transferAmount)

			},
		},
		{
			"Transfer with pre-set ParachainIBCTokenInfo",
			sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
			func() {
				// Add parachain token info
				chainBtransMiddleware := suite.chainB.TransferMiddleware()
				expChainBBalanceDiff = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
				err := chainBtransMiddleware.AddParachainIBCInfo(suite.chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
				suite.Require().NoError(err)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.SetupTest() // reset

			path = NewTransferPath(suite.chainA, suite.chainB)
			suite.coordinator.Setup(path)

			tc.malleate()

			originalChainABalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
			// chainB.SenderAccount: 10000000000000000000stake
			originalChainBBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())

			fmt.Println("chainB.AllBalances(chainB.SenderAccount.GetAddress())", suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress()))
			msg := ibctransfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, coinToSendToB, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
			_, err := suite.chainA.SendMsgs(msg)
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
			suite.Require().Equal(originalChainABalance.Sub(tc.expChainABalanceDiff), newChainABalance)

			// and dest chain balance contains voucher
			expBalance := originalChainBBalance.Add(expChainBBalanceDiff)
			gotBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())
			fmt.Println("expBalance", expBalance)
			fmt.Println("gotBalance", gotBalance)
			suite.Require().Equal(expBalance, gotBalance)
		})
	}
}

// TODO: use testsuite here.
func (suite *TransferMiddlewareTestSuite) TestSendPacket() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		// when transfer via sdk transfer from A (module) -> B (contract)
		nativeToken   = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
		timeoutHeight = clienttypes.NewHeight(1, 110)
	)
	var (
		expChainBBalanceDiff sdk.Coin
		path                 = NewTransferPath(suite.chainA, suite.chainB)
		expChainABalanceDiff = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
	)

	suite.SetupTest() // reset

	path = NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	// Add parachain token info
	chainBtransMiddlewareKeeper := suite.chainB.TransferMiddleware()
	expChainBBalanceDiff = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
	err := chainBtransMiddlewareKeeper.AddParachainIBCInfo(suite.chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
	suite.Require().NoError(err)

	originalChainABalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	originalChainBBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())

	msg := ibctransfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nativeToken, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
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
	expBalance := originalChainBBalance.Add(expChainBBalanceDiff)
	gotBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())
	suite.Require().Equal(expBalance, gotBalance)

	// send token back
	msg = ibctransfertypes.NewMsgTransfer(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, nativeToken, suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
	_, err = suite.chainB.SendMsgs(msg)
	suite.Require().NoError(err)
	suite.Require().NoError(err, path.EndpointA.UpdateClient())

	// then
	suite.Require().Equal(1, len(suite.chainB.PendingSendPackets))

	// and when relay to chain B and handle Ack on chain A
	err = suite.coordinator.RelayAndAckPendingPacketsReverse(path)
	suite.Require().NoError(err)

	// then
	suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

	// check escrow address don't have any token in chain B
	escrowAddressChainB := ibctransfertypes.GetEscrowAddress(ibctransfertypes.PortID, path.EndpointB.ChannelID)
	escrowTokenChainB := suite.chainB.AllBalances(escrowAddressChainB)
	suite.Require().Equal(sdk.Coins{}, escrowTokenChainB)

	// check escrow address don't have any token in chain A
	escrowAddressChainA := ibctransfertypes.GetEscrowAddress(ibctransfertypes.PortID, path.EndpointA.ChannelID)
	escrowTokenChainA := suite.chainA.AllBalances(escrowAddressChainA)
	suite.Require().Equal(sdk.Coins{}, escrowTokenChainA)

	// equal chain A sender address balances
	chainASenderBalances := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	suite.Require().Equal(originalChainABalance, chainASenderBalances)
}

// TODO: use testsuite here.
func (suite *TransferMiddlewareTestSuite) TestTimeOutPacket() {
	var (
		transferAmount = sdk.NewInt(1000000000)
		// when transfer via sdk transfer from A (module) -> B (contract)
		nativeToken   = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
		timeoutHeight = clienttypes.NewHeight(1, 110)
	)
	var (
		expChainBBalanceDiff sdk.Coin
		path                 = NewTransferPath(suite.chainA, suite.chainB)
		expChainABalanceDiff = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
	)

	suite.SetupTest() // reset

	path = NewTransferPath(suite.chainA, suite.chainB)
	suite.coordinator.Setup(path)

	// Add parachain token info
	chainBtransMiddlewareKeeper := suite.chainB.TransferMiddleware()
	expChainBBalanceDiff = sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)
	err := chainBtransMiddlewareKeeper.AddParachainIBCInfo(suite.chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
	suite.Require().NoError(err)

	originalChainABalance := suite.chainA.AllBalances(suite.chainA.SenderAccount.GetAddress())
	originalChainBBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())

	msg := ibctransfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, nativeToken, suite.chainA.SenderAccount.GetAddress().String(), suite.chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
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
	expBalance := originalChainBBalance.Add(expChainBBalanceDiff)
	gotBalance := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())
	suite.Require().Equal(expBalance, gotBalance)

	// send token back
	timeout := uint64(suite.chainB.LastHeader.Header.Time.Add(time.Nanosecond).UnixNano()) // will timeout
	msg = ibctransfertypes.NewMsgTransfer(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, nativeToken, suite.chainB.SenderAccount.GetAddress().String(), suite.chainA.SenderAccount.GetAddress().String(), clienttypes.NewHeight(1, 20), timeout, "")
	_, err = suite.chainB.SendMsgs(msg)
	suite.Require().NoError(err)
	suite.Require().NoError(err, path.EndpointA.UpdateClient())

	// then
	suite.Require().Equal(1, len(suite.chainB.PendingSendPackets))
	// and when relay to chain B and handle Ack on chain A
	err = suite.coordinator.TimeoutPendingPacketsReverse(path)
	suite.Require().NoError(err)

	// then
	suite.Require().Equal(0, len(suite.chainB.PendingSendPackets))

	// equal chain A sender address balances
	chainBSenderBalances := suite.chainB.AllBalances(suite.chainB.SenderAccount.GetAddress())
	suite.Equal(expBalance, chainBSenderBalances)
}
