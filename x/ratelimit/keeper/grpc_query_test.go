package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

func (s *KeeperTestSuite) TestGRPCAllRateLimits() {
	var (
		denom              = sdk.DefaultBondDenom // Should use DefaultBondDenom; otherwise it returns channel value error
		channelID          = "channel-0"
		maxPercentSend     = sdkmath.NewInt(50)
		maxPercentRecv     = sdkmath.NewInt(50)
		minRateLimitAmount = sdkmath.NewInt(100000)
		durationHours      = uint64(1)
	)
	s.addRateLimit(denom, channelID, maxPercentSend, maxPercentRecv, minRateLimitAmount, durationHours)

	for _, tc := range []struct {
		name      string
		req       *types.QueryAllRateLimitsRequest
		expectErr bool
		postRun   func(*types.QueryAllRateLimitsResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"happy case",
			&types.QueryAllRateLimitsRequest{},
			false,
			func(resp *types.QueryAllRateLimitsResponse) {
				s.Require().Len(resp.GetRateLimits(), 1)
				s.Require().Equal(denom, resp.RateLimits[0].Path.Denom)
				s.Require().Equal(channelID, resp.RateLimits[0].Path.ChannelID)
				s.Require().Equal(maxPercentSend, resp.RateLimits[0].Quota.MaxPercentSend)
				s.Require().Equal(maxPercentRecv, resp.RateLimits[0].Quota.MaxPercentRecv)
				s.Require().Equal(minRateLimitAmount, resp.RateLimits[0].MinRateLimitAmount)
				s.Require().Equal(durationHours, resp.RateLimits[0].Quota.DurationHours)
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.AllRateLimits(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}
