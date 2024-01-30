package types

import (
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

func TestAnyPackUnpack(t *testing.T) {
	registry := types.NewInterfaceRegistry()
	registry.RegisterInterface("centauri.ratelimit.v1beta1.MsgAddRateLimit", (*sdk.Msg)(nil))
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddRateLimitLegacy{},
	)

	input := &MsgAddRateLimitLegacy{
		Authority:          authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		Denom:              "test",
		ChannelId:          "test",
		MaxPercentSend:     math.NewInt(10),
		MaxPercentRecv:     math.NewInt(10),
		DurationHours:      1000,
		MinRateLimitAmount: math.NewInt(10),
	}
	var msg sdk.Msg

	// with cache
	any, err := types.NewAnyWithValue(input)
	require.NoError(t, err)
	require.Equal(t, input, any.GetCachedValue())
	err = registry.UnpackAny(any, &msg)

	require.NoError(t, err)
	require.Equal(t, input, msg)
}
