package keeper_test

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	// ibctesting "github.com/cosmos/ibc-go/v7/testing"

	ibctesting "github.com/notional-labs/centauri/v5/app/ibctesting"
	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

func (s *KeeperTestSuite) TestGRPCAllRateLimits() {
	// Add some sample rate limits
	s.SetupSampleRateLimits(sampleRateLimitA, sampleRateLimitB, sampleRateLimitC, sampleRateLimitD)

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
				s.Require().Len(resp.GetRateLimits(), 4)
				s.Require().Equal(sampleRateLimitA.Denom, resp.RateLimits[0].Path.Denom)
				s.Require().Equal(sampleRateLimitA.ChannelID, resp.RateLimits[0].Path.ChannelID)
				s.Require().Equal(sampleRateLimitA.MaxPercentSend, resp.RateLimits[0].Quota.MaxPercentSend)
				s.Require().Equal(sampleRateLimitA.MaxPercentRecv, resp.RateLimits[0].Quota.MaxPercentRecv)
				s.Require().Equal(sampleRateLimitA.MinRateLimitAmount, resp.RateLimits[0].MinRateLimitAmount)
				s.Require().Equal(sampleRateLimitA.DurationHours, resp.RateLimits[0].Quota.DurationHours)
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

func (s *KeeperTestSuite) TestGRPCRateLimit() {
	// Add some sample rate limits
	s.SetupSampleRateLimits(sampleRateLimitA, sampleRateLimitB, sampleRateLimitC, sampleRateLimitD)

	for _, tc := range []struct {
		name      string
		req       *types.QueryRateLimitRequest
		expectErr bool
		postRun   func(*types.QueryRateLimitResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"happy case",
			&types.QueryRateLimitRequest{
				Denom:     sampleRateLimitA.Denom,
				ChannelID: sampleRateLimitA.ChannelID,
			},
			false,
			func(resp *types.QueryRateLimitResponse) {
				s.Require().Equal(sampleRateLimitA.Denom, resp.RateLimit.Path.Denom)
				s.Require().Equal(sampleRateLimitA.ChannelID, resp.RateLimit.Path.ChannelID)
				s.Require().Equal(sampleRateLimitA.MaxPercentSend, resp.RateLimit.Quota.MaxPercentSend)
				s.Require().Equal(sampleRateLimitA.MaxPercentRecv, resp.RateLimit.Quota.MaxPercentRecv)
				s.Require().Equal(sampleRateLimitA.MinRateLimitAmount, resp.RateLimit.MinRateLimitAmount)
				s.Require().Equal(sampleRateLimitA.DurationHours, resp.RateLimit.Quota.DurationHours)
			},
		},
		{
			"query by invalid denom",
			&types.QueryRateLimitRequest{
				Denom:     "invalidDenom",
				ChannelID: sampleRateLimitA.ChannelID,
			},
			true,
			nil,
		},
		{
			"query by invalid channel id",
			&types.QueryRateLimitRequest{
				Denom:     sampleRateLimitA.Denom,
				ChannelID: "invalidChannelID",
			},
			true,
			nil,
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.RateLimit(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRateLimitsByChainID() {
	// Add some sample rate limits
	s.SetupSampleRateLimits(sampleRateLimitA, sampleRateLimitB)

	// // Create client and connections on both chains
	// path := ibctesting.NewPath(s.chainA, s.chainB)
	// s.coordinator.SetupConnections(path)
	// path.SetChannelOrdered()

	// // Initialize channel
	// err := path.EndpointA.ChanOpenInit()
	// s.Require().NoError(err)

	for _, tc := range []struct {
		name      string
		req       *types.QueryRateLimitsByChainIDRequest
		expectErr bool
		postRun   func(*types.QueryRateLimitsByChainIDResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"happy case",
			&types.QueryRateLimitsByChainIDRequest{
				ChainId: s.chainA.ChainID,
			},
			false,
			func(resp *types.QueryRateLimitsByChainIDResponse) {
				fmt.Println("resp: ", resp)
				// s.Require().Len(resp.GetRateLimits(), 2)
			},
		},
	} {
		s.Run(tc.name, func() {
			ctx := sdk.WrapSDKContext(s.chainA.GetContext())

			resp, err := s.querier.RateLimitsByChainID(ctx, tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestRateLimitsByChannelID() {
	// Add some sample rate limits
	s.SetupSampleRateLimits(sampleRateLimitA, sampleRateLimitB, sampleRateLimitC, sampleRateLimitD)

	for _, tc := range []struct {
		name      string
		req       *types.QueryRateLimitsByChannelIDRequest
		expectErr bool
		postRun   func(*types.QueryRateLimitsByChannelIDResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"happy case",
			&types.QueryRateLimitsByChannelIDRequest{
				ChannelID: sampleRateLimitA.ChannelID,
			},
			false,
			func(resp *types.QueryRateLimitsByChannelIDResponse) {
				s.Require().Len(resp.GetRateLimits(), 2)
				s.Require().Equal(sampleRateLimitA.Denom, resp.RateLimits[0].Path.Denom)
				s.Require().Equal(sampleRateLimitA.ChannelID, resp.RateLimits[0].Path.ChannelID)
				s.Require().Equal(sampleRateLimitA.MaxPercentSend, resp.RateLimits[0].Quota.MaxPercentSend)
				s.Require().Equal(sampleRateLimitA.MaxPercentRecv, resp.RateLimits[0].Quota.MaxPercentRecv)
				s.Require().Equal(sampleRateLimitA.MinRateLimitAmount, resp.RateLimits[0].MinRateLimitAmount)
				s.Require().Equal(sampleRateLimitA.DurationHours, resp.RateLimits[0].Quota.DurationHours)
				s.Require().Equal(sampleRateLimitB.Denom, resp.RateLimits[1].Path.Denom)
				s.Require().Equal(sampleRateLimitB.ChannelID, resp.RateLimits[1].Path.ChannelID)
				s.Require().Equal(sampleRateLimitB.MaxPercentSend, resp.RateLimits[1].Quota.MaxPercentSend)
				s.Require().Equal(sampleRateLimitB.MaxPercentRecv, resp.RateLimits[1].Quota.MaxPercentRecv)
				s.Require().Equal(sampleRateLimitB.MinRateLimitAmount, resp.RateLimits[1].MinRateLimitAmount)
				s.Require().Equal(sampleRateLimitB.DurationHours, resp.RateLimits[1].Quota.DurationHours)
			},
		},
		{
			"query by chain id that does not exist",
			&types.QueryRateLimitsByChannelIDRequest{
				ChannelID: "channel-10",
			},
			false,
			func(resp *types.QueryRateLimitsByChannelIDResponse) {
				s.Require().Empty(resp.RateLimits)
			},
		},
		{
			"query by invalid chain id",
			&types.QueryRateLimitsByChannelIDRequest{
				ChannelID: "invalid/ChannelID",
			},
			true,
			nil,
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.RateLimitsByChannelID(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestAllWhitelistedAddresses() {
	// Add some sample whitelisted addresses
	whitelistedAddrPairs := []types.WhitelistedAddressPair{
		{Sender: s.addr(1).String(), Receiver: s.addr(2).String()},
		{Sender: s.addr(3).String(), Receiver: s.addr(4).String()},
		{Sender: s.addr(5).String(), Receiver: s.addr(6).String()},
	}

	for _, wap := range whitelistedAddrPairs {
		s.keeper.SetWhitelistedAddressPair(s.ctx, wap)
	}

	for _, tc := range []struct {
		name      string
		req       *types.QueryAllWhitelistedAddressesRequest
		expectErr bool
		postRun   func(*types.QueryAllWhitelistedAddressesResponse)
	}{
		{
			"nil request",
			nil,
			true,
			nil,
		},
		{
			"happy case",
			&types.QueryAllWhitelistedAddressesRequest{},
			false,
			func(resp *types.QueryAllWhitelistedAddressesResponse) {
				s.Require().Len(resp.GetAddressPairs(), 3)
			},
		},
	} {
		s.Run(tc.name, func() {
			resp, err := s.querier.AllWhitelistedAddresses(sdk.WrapSDKContext(s.ctx), tc.req)
			if tc.expectErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
				tc.postRun(resp)
			}
		})
	}
}

func (s *KeeperTestSuite) TestChanOpenInit() {
	// Create client and connections on both chains
	path := ibctesting.NewPath(s.chainA, s.chainB)
	s.coordinator.SetupConnections(path)

	path.SetChannelOrdered()

	// Initialize channel
	//
	// Works well with ibctesting "github.com/cosmos/ibc-go/v7/testing" but it does not working due to the following error message with customibctesting
	// Error: could not retrieve module from port-id: ports/mock: capability not found
	//
	err := path.EndpointA.ChanOpenInit()
	s.Require().NoError(err)

	storedChannel, found := s.chainA.App.GetIBCKeeper().ChannelKeeper.GetChannel(s.chainA.GetContext(), path.EndpointA.ChannelConfig.PortID, path.EndpointA.ChannelID)
	s.True(found)
	fmt.Println("storedChannel: ", storedChannel)
}
