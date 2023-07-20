package keeper_test

import (
	"testing"

	helpers "github.com/notional-labs/centauri/v3/app/helpers"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	"github.com/stretchr/testify/require"
)

func TestParaTokenInfo(t *testing.T) {
	app := helpers.SetupCentauriAppWithValSet(t)
	ctx := helpers.NewContextForApp(*app)

	tokenInfos := make([]types.ParachainIBCTokenInfo, 1)
	tokenInfos[0] = types.ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test",
		ChannelId:   "channel-0",
		NativeDenom: "pica",
		AssetId:     "1",
	}
	app.TransferMiddlewareKeeper.InitGenesis(ctx, types.GenesisState{
		TokenInfos: tokenInfos,
	})

	info, err := app.TransferMiddlewareKeeper.ParaTokenInfo(ctx, &types.QueryParaTokenInfoRequest{NativeDenom: "pica"})

	require.NoError(t, err)
	require.Equal(t, "1", info.AssetId)
	require.Equal(t, "pica", info.NativeDenom)
	require.Equal(t, "ibc-test", info.IbcDenom)
	require.Equal(t, "channel-0", info.ChannelId)
}

func TestEscrowAddress(t *testing.T) {
	app := helpers.SetupCentauriAppWithValSet(t)
	ctx := helpers.NewContextForApp(*app)

	tokenInfos := make([]types.ParachainIBCTokenInfo, 1)
	tokenInfos[0] = types.ParachainIBCTokenInfo{
		IbcDenom:    "ibc-test",
		ChannelId:   "channel-0",
		NativeDenom: "pica",
		AssetId:     "1",
	}
	app.TransferMiddlewareKeeper.InitGenesis(ctx, types.GenesisState{
		TokenInfos: tokenInfos,
	})

	escrowResponse, err := app.TransferMiddlewareKeeper.EscrowAddress(ctx, &types.QueryEscrowAddressRequest{ChannelId: "channel-0"})
	require.NoError(t, err)
	require.Equal(t, escrowResponse.EscrowAddress, transfertypes.GetEscrowAddress(transfertypes.PortID, "channel-0").String())
}