package keeper_test

import (
	"testing"
	"github.com/notional-labs/centauri/v3/app/helpers"
	"github.com/notional-labs/centauri/v3/x/transfermiddleware/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	app := helpers.SetupCentauriAppWithValSet(t)
	ctx := helpers.NewContextForApp(*app)
	params := types.DefaultParams()

 	app.TransferMiddlewareKeeper.SetParams(ctx, params)

	require.Equal(t, params, app.TransferMiddlewareKeeper.GetParams(ctx))
}
