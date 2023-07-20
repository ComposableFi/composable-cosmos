package keeper_test

import (
	"testing"

	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"

	"github.com/stretchr/testify/require"
)

func TestBeginBlocker(t *testing.T) {
	app, ctx := SetupTest(t)

	infos := [5]types.ParachainIBCTokenInfo{
		{
			IbcDenom:    "ibc-test-1",
			ChannelId:   "channel-1",
			NativeDenom: "native-1",
			AssetId:     "2",
		},
		{
			IbcDenom:    "ibc-test-2",
			ChannelId:   "channel-1",
			NativeDenom: "native-2",
			AssetId:     "3",
		},
		{
			IbcDenom:    "ibc-test-3",
			ChannelId:   "channel-1",
			NativeDenom: "native-3",
			AssetId:     "4",
		},
		{
			IbcDenom:    "ibc-test-4",
			ChannelId:   "channel-2",
			NativeDenom: "native-4",
			AssetId:     "5",
		},
		{
			IbcDenom:    "ibc-test-5",
			ChannelId:   "channel-2",
			NativeDenom: "native-5",
			AssetId:     "6",
		},
	}

	count := 0
	for _, info := range infos {
		app.TransferMiddlewareKeeper.AddParachainIBCInfo(
			ctx, 
			info.IbcDenom,
			info.ChannelId,
			info.NativeDenom,
			info.AssetId,
		)
	}

	app.TransferMiddlewareKeeper.IterateRemoveListInfo(ctx, func(_ types.RemoveParachainIBCTokenInfo) (stop bool) {
		count++
		return false
	})

	require.Equal(t, count, 0)
	
	for _, info := range infos {
		app.TransferMiddlewareKeeper.AddParachainIBCInfoToRemoveList(
			ctx, 
			info.NativeDenom,
		)
	}

	app.TransferMiddlewareKeeper.IterateRemoveListInfo(ctx, func(_ types.RemoveParachainIBCTokenInfo) (stop bool) {
		count++
		return false
	})

	require.Equal(t, count, 5)

	app.TransferMiddlewareKeeper.BeginBlocker(ctx)

	countRemove := 0
	app.TransferMiddlewareKeeper.IterateRemoveListInfo(ctx, func(removeList types.RemoveParachainIBCTokenInfo) (stop bool) {
		if ctx.BlockTime().After(removeList.RemoveTime) {
			t.Log(ctx.BlockTime(), removeList.RemoveTime)
			countRemove++
			count--
		}
		return false
	})

	require.Equal(t, countRemove + count, 5)
}