package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingmiddlewaretypes "github.com/notional-labs/composable/v6/x/stakingmiddleware/types"
)

func (s *KeeperTestSuite) TestMsgUpdateEpochParams() {
	ctx, msgServer := s.ctx, s.msgServer
	require := s.Require()
	testCases := []struct {
		name      string
		input     *stakingmiddlewaretypes.MsgUpdateEpochParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "valid params",
			input: &stakingmiddlewaretypes.MsgUpdateEpochParams{
				Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
				Params: stakingmiddlewaretypes.Params{
					BlocksPerEpoch:                           20,
					AllowUnbondAfterEpochProgressBlockNumber: 7,
				},
			},
			expErr: false,
		},
		{
			name: "invalid authority",
			input: &stakingmiddlewaretypes.MsgUpdateEpochParams{
				Authority: "invalid",
				Params: stakingmiddlewaretypes.Params{
					BlocksPerEpoch:                           20,
					AllowUnbondAfterEpochProgressBlockNumber: 7,
				},
			},
			expErr:    true,
			expErrMsg: "invalid authority",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.UpdateEpochParams(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
				require.Contains(err.Error(), tc.expErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAddRevenueFundsToStaking() {
	accAddrs := []sdk.AccAddress{
		sdk.AccAddress([]byte("addr1_______________")),
	}
	ctx, msgServer := s.ctx, s.msgServer
	require := s.Require()
	feeCoin := sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(150))
	feeAmount := sdk.NewCoins(feeCoin)
	s.app.BankKeeper.MintCoins(ctx, minttypes.ModuleName, feeAmount)
	s.app.BankKeeper.SendCoinsFromModuleToAccount(ctx, minttypes.ModuleName, accAddrs[0], feeAmount)
	testCases := []struct {
		name      string
		input     *stakingmiddlewaretypes.MsgAddRevenueFundsToStakingParams
		expErr    bool
		expErrMsg string
	}{
		{
			name: "success",
			input: &stakingmiddlewaretypes.MsgAddRevenueFundsToStakingParams{
				FromAddress: accAddrs[0].String(),
				Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1))),
			},
			expErr: false,
		},
		{
			name: "invalid coin",
			input: &stakingmiddlewaretypes.MsgAddRevenueFundsToStakingParams{
				FromAddress: accAddrs[0].String(),
				Amount:      sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1))),
			},
			expErr:    true,
			expErrMsg: "Invalid coin",
		},
		{
			name: "invalid coin: two different coin denom",
			input: &stakingmiddlewaretypes.MsgAddRevenueFundsToStakingParams{
				FromAddress: accAddrs[0].String(),
				Amount:      sdk.NewCoins(sdk.NewCoin(sdk.DefaultBondDenom, sdk.NewInt(1)), sdk.NewCoin("test", sdk.NewInt(1))),
			},
			expErr:    true,
			expErrMsg: "Invalid coin",
		},
	}

	for _, tc := range testCases {
		tc := tc
		s.T().Run(tc.name, func(t *testing.T) {
			_, err := msgServer.AddRevenueFundsToStaking(ctx, tc.input)
			if tc.expErr {
				require.Error(err)
				require.Contains(err.Error(), tc.expErrMsg)
			} else {
				require.NoError(err)
			}
		})
	}
}
