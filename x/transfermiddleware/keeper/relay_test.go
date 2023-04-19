package keeper_test

import (
	"encoding/json"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v7/modules/core/04-channel/types"
	"github.com/golang/mock/gomock"
	"github.com/notional-labs/banksy/v2/test"
	"github.com/stretchr/testify/require"
)

var (
	testDenom  = "uatom"
	testAmount = "100"

	testSourcePort         = "transfer"
	testSourceChannel      = "channel-10"
	testDestinationPort    = "transfer"
	testDestinationChannel = "channel-11"

	sender   = "cosmos1vzxkv3lxccnttr9rs0002s93sgw72h7ghukuhs"
	receiver = "cosmos1q4p4gx889lfek5augdurrjclwtqvjhuntm6j4m"
)

func emptyPacket() channeltypes.Packet {
	return channeltypes.Packet{}
}

func transferPacket(t *testing.T, sender string, receiver string, metadata any) channeltypes.Packet {
	t.Helper()
	transferPacket := transfertypes.FungibleTokenPacketData{
		Denom:    testDenom,
		Amount:   testAmount,
		Sender:   sender,
		Receiver: receiver,
	}

	if metadata != nil {
		if mStr, ok := metadata.(string); ok {
			transferPacket.Memo = mStr
		} else {
			memo, err := json.Marshal(metadata)
			require.NoError(t, err)
			transferPacket.Memo = string(memo)
		}
	}

	transferData, err := transfertypes.ModuleCdc.MarshalJSON(&transferPacket)
	require.NoError(t, err)

	return channeltypes.Packet{
		SourcePort:         testSourcePort,
		SourceChannel:      testSourceChannel,
		DestinationPort:    testDestinationPort,
		DestinationChannel: testDestinationChannel,
		Data:               transferData,
	}
}

func TestOnRecvPacket_EmptyPacket(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	setup := test.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	cdc := setup.Initializer.Marshaler
	ibcMiddleware := setup.IBCMiddleware

	// Test data
	senderAccAddr := test.AccAddress()
	packet := emptyPacket()

	ack := ibcMiddleware.OnRecvPacket(ctx, packet, senderAccAddr)
	require.False(t, ack.Success())

	expectedAck := &channeltypes.Acknowledgement{}
	err := cdc.UnmarshalJSON(ack.Acknowledgement(), expectedAck)
	require.NoError(t, err)
	require.Equal(t, "ABCI code: 1: error handling packet: see events for details", expectedAck.GetError())
}

func TestOnRecvPacket_TransferPacket(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	setup := test.NewTestSetup(t, ctl)
	ctx := setup.Initializer.Ctx
	cdc := setup.Initializer.Marshaler
	ibcMiddleware := setup.IBCMiddleware

	// Test data
	senderAccAddr := test.AccAddress()
	packet := transferPacket(t, sender, receiver, nil)
	setup.Mocks.TransferKeeperMock.EXPECT().GetReceiveEnabled(ctx).Return(true).AnyTimes()
	prefixedDenom := "transfer/channel-11/uatom"
	denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)
	traceHash := denomTrace.Hash()
	setup.Mocks.TransferKeeperMock.EXPECT().HasDenomTrace(ctx, traceHash).Return(true).AnyTimes()
	amount, _ := sdk.NewIntFromString("100")
	voucher := sdk.NewCoin("ibc/5FEB332D2B121921C792F1A0DBF7C3163FF205337B4AFE6E14F69E8E49545F49", amount)
	setup.Mocks.BankKeeperMock.EXPECT().MintCoins(ctx, "transfermiddleware", sdk.NewCoins(voucher)).Return(nil).AnyTimes()
	receiver, _ := sdk.AccAddressFromBech32(receiver)
	setup.Mocks.BankKeeperMock.EXPECT().SendCoinsFromModuleToAccount(ctx, "transfermiddleware", receiver, sdk.NewCoins(voucher))
	ack := ibcMiddleware.OnRecvPacket(ctx, packet, senderAccAddr)
	fmt.Println("ack", string(ack.Acknowledgement()))
	require.True(t, ack.Success())

	expectedAck := &channeltypes.Acknowledgement{}
	err := cdc.UnmarshalJSON(ack.Acknowledgement(), expectedAck)
	require.NoError(t, err)
	require.Equal(t, "", expectedAck.GetError())

	require.Equal(t, []byte{0x1}, expectedAck.GetResult())
}
