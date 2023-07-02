package keeper

import (
	"context"

	"github.com/notional-labs/centauri/v3/x/ratelimit/types"
)

var _ types.QueryServer = Keeper{}

// Query all rate limits
func (k Keeper) AllRateLimits(goCtx context.Context, req *types.QueryAllRateLimitsRequest) (*types.QueryAllRateLimitsResponse, error) {
}

// Query a rate limit by denom and channelId
func (k Keeper) RateLimit(goCtx context.Context, req *types.QueryRateLimitRequest) (*types.QueryRateLimitResponse, error) {
}

// Query all rate limits for a given chain
func (k Keeper) RateLimitsByChainId(c context.Context, req *types.QueryRateLimitsByChainIdRequest) (*types.QueryRateLimitsByChainIdResponse, error) {
}

// Query all rate limits for a given channel
func (k Keeper) RateLimitsByChannelId(c context.Context, req *types.QueryRateLimitsByChannelIdRequest) (*types.QueryRateLimitsByChannelIdResponse, error) {
}

// Query all blacklisted denoms
func (k Keeper) AllBlacklistedDenoms(c context.Context, req *types.QueryAllBlacklistedDenomsRequest) (*types.QueryAllBlacklistedDenomsResponse, error) {
}

// Query all whitelisted addresses
func (k Keeper) AllWhitelistedAddresses(c context.Context, req *types.QueryAllWhitelistedAddressesRequest) (*types.QueryAllWhitelistedAddressesResponse, error) {
}
