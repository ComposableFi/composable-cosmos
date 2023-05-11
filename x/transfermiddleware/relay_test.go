package transfermiddleware_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	clienttypes "github.com/cosmos/ibc-go/v7/modules/core/02-client/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	ibctesting "github.com/cosmos/ibc-go/v7/testing"
	customibctesting "github.com/notional-labs/banksy/v2/app/ibctesting"
	routertypes "github.com/strangelove-ventures/packet-forward-middleware/v7/router/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO: use testsuite here.
func TestOnrecvPacket(t *testing.T) {
	// scenario: given two chains,
	//           with a contract on chain B
	//           then the contract can handle the receiving side of an ics20 transfer
	//           that was started on chain A via ibc transfer module

	transferAmount := sdk.NewInt(1000000000)
	var (
		coordinator = customibctesting.NewCoordinator(t, 2)
		chainA      = coordinator.GetChain(customibctesting.GetChainID(0))
		chainB      = coordinator.GetChain(customibctesting.GetChainID(1))
	)
	coordinator.CommitBlock(chainA, chainB)

	path := customibctesting.NewPath(chainA, chainB)
	path.EndpointA.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  "transfer",
		Version: ibctransfertypes.Version,
		Order:   channeltypes.UNORDERED,
	}
	path.EndpointB.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  "transfer",
		Version: ibctransfertypes.Version,
		Order:   channeltypes.UNORDERED,
	}

	coordinator.SetupConnections(path)
	coordinator.CreateChannels(path)

	// when transfer via sdk transfer from A (module) -> B (contract)
	coinToSendToB := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1000000000))
	timeoutHeight := clienttypes.NewHeight(1, 110)

	testCases := []struct {
		name                 string
		expChainABalanceDiff sdk.Coin
		expChainBBalanceDiff sdk.Coin
		malleate             func()
	}{
		{
			"Transfer with no pre-set ParachainIBCTokenInfo",
			sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
			ibctransfertypes.GetTransferCoin(path.EndpointB.ChannelConfig.PortID, path.EndpointB.ChannelID, coinToSendToB.Denom, transferAmount),
			func() {},
		},
		{
			"Transfer with pre-set ParachainIBCTokenInfo",
			sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
			sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
			func() {
				// Add parachain token info
				chainBtransMiddleware := chainB.TransferMiddleware()
				err := chainBtransMiddleware.AddParachainIBCInfo(chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
				require.NoError(t, err)
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			tc.malleate()

			originalChainABalance := chainA.AllBalances(chainA.SenderAccount.GetAddress())
			// chainB.SenderAccount: 10000000000000000000stake
			originalChainBBalance := chainB.AllBalances(chainB.SenderAccount.GetAddress())

			fmt.Println("chainB.AllBalances(chainB.SenderAccount.GetAddress())", chainB.AllBalances(chainB.SenderAccount.GetAddress()))
			msg := ibctransfertypes.NewMsgTransfer(path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID, coinToSendToB, chainA.SenderAccount.GetAddress().String(), chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, "")
			_, err := chainA.SendMsgs(msg)
			require.NoError(t, err)
			require.NoError(t, path.EndpointB.UpdateClient())

			// then
			require.Equal(t, 1, len(chainA.PendingSendPackets))
			require.Equal(t, 0, len(chainB.PendingSendPackets))

			// and when relay to chain B and handle Ack on chain A
			err = coordinator.RelayAndAckPendingPackets(path)
			require.NoError(t, err)

			// then
			require.Equal(t, 0, len(chainA.PendingSendPackets))
			require.Equal(t, 0, len(chainB.PendingSendPackets))

			// and source chain balance was decreased
			newChainABalance := chainA.AllBalances(chainA.SenderAccount.GetAddress())
			assert.Equal(t, originalChainABalance.Sub(tc.expChainABalanceDiff), newChainABalance)

			// and dest chain balance contains voucher
			expBalance := originalChainBBalance.Add(tc.expChainBBalanceDiff)
			gotBalance := chainB.AllBalances(chainB.SenderAccount.GetAddress())
			fmt.Println("expBalance", expBalance)
			fmt.Println("gotBalance", gotBalance)
			assert.Equal(t, expBalance, gotBalance)
		})
	}
}

func TestOnrecvPacket2(t *testing.T) {
	// scenario: given two chains,
	//           with a contract on chain B
	//           then the contract can handle the receiving side of an ics20 transfer
	//           that was started on chain A via ibc transfer module

	transferAmount := sdk.NewInt(1)
	var (
		coordinator = customibctesting.NewCoordinator(t, 3)
		chainA      = coordinator.GetChain(customibctesting.GetChainID(0))
		chainB      = coordinator.GetChain(customibctesting.GetChainID(1))
		chainC      = coordinator.GetChain(customibctesting.GetChainID(2))
	)
	coordinator.CommitBlock(chainA, chainB, chainC)

	pathAB := customibctesting.NewPath(chainA, chainB)
	pathBC := customibctesting.NewPath(chainB, chainC)

	pathAB.EndpointA.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  "transfer",
		Version: ibctransfertypes.Version,
		Order:   channeltypes.UNORDERED,
	}
	pathAB.EndpointB.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  "transfer",
		Version: ibctransfertypes.Version,
		Order:   channeltypes.UNORDERED,
	}

	pathBC.EndpointA.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  "transfer",
		Version: ibctransfertypes.Version,
		Order:   channeltypes.UNORDERED,
	}
	pathBC.EndpointB.ChannelConfig = &ibctesting.ChannelConfig{
		PortID:  "transfer",
		Version: ibctransfertypes.Version,
		Order:   channeltypes.UNORDERED,
	}
	coordinator.SetupConnections(pathAB)
	coordinator.CreateChannels(pathAB)
	coordinator.SetupConnections(pathBC)
	coordinator.CreateChannels(pathBC)
	// when transfer via sdk transfer from A -> B -> C
	coinASendToB := sdk.NewCoin(sdk.DefaultBondDenom, transferAmount)

	timeoutHeight := clienttypes.NewHeight(1, 110)

	ibcDenomAtoB := ibctransfertypes.GetPrefixedDenom(pathAB.EndpointB.ChannelConfig.PortID, pathAB.EndpointB.ChannelID, sdk.DefaultBondDenom)
	testCases := []struct {
		name                 string
		expChainABalanceDiff sdk.Coin
		expChainBBalanceDiff sdk.Coin
		expChainCBalanceDiff sdk.Coin
		malleate             func()
	}{
		{
			name:                 "Transfer with no pre-set ParachainIBCTokenInfo",
			expChainABalanceDiff: sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
			expChainCBalanceDiff: ibctransfertypes.GetTransferCoin(pathBC.EndpointB.ChannelConfig.PortID, pathBC.EndpointB.ChannelID, ibcDenomAtoB, transferAmount),
			malleate:             func() {},
		},
		// {
		// 	"Transfer with pre-set ParachainIBCTokenInfo",
		// 	sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(2000000000)),
		// 	sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
		// 	sdk.NewCoin(sdk.DefaultBondDenom, transferAmount),
		// 	func() {
		// 		// Add parachain token info
		// 		chainBtransMiddleware := chainB.TransferMiddleware()
		// 		err := chainBtransMiddleware.AddParachainIBCInfo(chainB.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
		// 		require.NoError(t, err)

		// 		chainCtransMiddleware := chainC.TransferMiddleware()
		// 		err = chainCtransMiddleware.AddParachainIBCInfo(chainC.GetContext(), "ibc/C053D637CCA2A2BA030E2C5EE1B28A16F71CCB0E45E8BE52766DC1B241B77878", "channel-0", sdk.DefaultBondDenom)
		// 		require.NoError(t, err)
		// 	},
		// },
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {

			tc.malleate()

			originalChainABalance := chainA.AllBalances(chainA.SenderAccount.GetAddress())
			// chainB.SenderAccount: 10000000000000000000stake
			// originalChainBBalance := chainB.AllBalances(chainB.SenderAccount.GetAddress())

			originalChainCBalance := chainC.AllBalances(chainC.SenderAccount.GetAddress())

			fmt.Println("Begin")
			fmt.Println("chainA.AllBalances(chainA.SenderAccount.GetAddress())", chainA.AllBalances(chainA.SenderAccount.GetAddress()))
			fmt.Println("chainB.AllBalances(chainB.SenderAccount.GetAddress())", chainB.AllBalances(chainB.SenderAccount.GetAddress()))
			fmt.Println("chainC.AllBalances(chainC.SenderAccount.GetAddress())", chainC.AllBalances(chainC.SenderAccount.GetAddress()))
			forwardMetadata := routertypes.PacketMetadata{
				Forward: &routertypes.ForwardMetadata{
					Receiver: chainC.SenderAccount.GetAddress().String(),
					Port:     "transfer",
					Channel:  pathBC.EndpointA.ChannelID,
				},
			}
			memo, err := json.Marshal(forwardMetadata)

			msg := ibctransfertypes.NewMsgTransfer(pathAB.EndpointA.ChannelConfig.PortID, pathAB.EndpointA.ChannelID, coinASendToB, chainA.SenderAccount.GetAddress().String(), chainB.SenderAccount.GetAddress().String(), timeoutHeight, 0, string(memo))
			_, err = chainA.SendMsgs(msg)
			require.NoError(t, err)
			require.NoError(t, pathAB.EndpointB.UpdateClient())
			require.NoError(t, pathBC.EndpointB.UpdateClient())

			// then
			require.Equal(t, 1, len(chainA.PendingSendPackets))
			require.Equal(t, 0, len(chainB.PendingSendPackets))

			// and when relay to chain B and handle Ack on chain A
			err = coordinator.RelayAndAckPendingPackets(pathAB)
			require.NoError(t, err)

			err = coordinator.RelayAndAckPendingPackets(pathBC)
			require.NoError(t, err)
			// then
			require.Equal(t, 0, len(chainA.PendingSendPackets))
			require.Equal(t, 0, len(chainB.PendingSendPackets))
			fmt.Println("After A -> B")
			fmt.Println("chainA.AllBalances(chainA.SenderAccount.GetAddress())", chainA.AllBalances(chainA.SenderAccount.GetAddress()))
			fmt.Println("chainB.AllBalances(chainB.SenderAccount.GetAddress())", chainB.AllBalances(chainB.SenderAccount.GetAddress()))
			fmt.Println("chainC.AllBalances(chainC.SenderAccount.GetAddress())", chainC.AllBalances(chainC.SenderAccount.GetAddress()))
			// and source chain balance was decreased
			newChainABalance := chainA.AllBalances(chainA.SenderAccount.GetAddress())
			assert.Equal(t, originalChainABalance.Sub(tc.expChainABalanceDiff), newChainABalance)

			// and dest chain balance contains voucher
			expBalance := originalChainCBalance.Add(tc.expChainCBalanceDiff)
			gotBalance := chainC.AllBalances(chainC.SenderAccount.GetAddress())
			fmt.Println("expBalance", expBalance)
			fmt.Println("gotBalance", gotBalance)
			assert.Equal(t, expBalance, gotBalance)
		})
	}
}
