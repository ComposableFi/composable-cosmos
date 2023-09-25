package keeper

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	transfertypes "github.com/cosmos/ibc-go/v7/modules/apps/transfer/types"
	ibctmtypes "github.com/cosmos/ibc-go/v7/modules/light-clients/07-tendermint"

	"github.com/notional-labs/centauri/v5/x/ratelimit/types"
)

var _ types.QueryServer = queryServer{}

type queryServer struct {
	Keeper
}

// NewQueryServer returns an implementation of the QueryServer
// for the provided Keeper.
func NewQueryServer(k Keeper) types.QueryServer {
	return queryServer{Keeper: k}
}

// AllRateLimits queries all rate limits.
func (q queryServer) AllRateLimits(c context.Context, req *types.QueryAllRateLimitsRequest) (*types.QueryAllRateLimitsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "empty request")
	}

	ctx := sdk.UnwrapSDKContext(c)
	rateLimits := q.GetAllRateLimits(ctx)
	return &types.QueryAllRateLimitsResponse{RateLimits: rateLimits}, nil
}

// RateLimit queries a rate limit by denom and channel id.
func (q queryServer) RateLimit(c context.Context, req *types.QueryRateLimitRequest) (*types.QueryRateLimitResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	rateLimit, found := q.GetRateLimit(ctx, req.Denom, req.ChannelID)
	if !found {
		return &types.QueryRateLimitResponse{}, nil
	}
	return &types.QueryRateLimitResponse{RateLimit: &rateLimit}, nil
}

// RateLimitsByChainID queries all rate limits for a given chain.
func (q queryServer) RateLimitsByChainID(c context.Context, req *types.QueryRateLimitsByChainIDRequest) (*types.QueryRateLimitsByChainIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	rateLimits := []types.RateLimit{}
	for _, rateLimit := range q.GetAllRateLimits(ctx) {

		// Determine the client state from the channel Id
		_, clientState, err := q.channelKeeper.GetChannelClientState(ctx, transfertypes.PortID, rateLimit.Path.ChannelID)
		if err != nil {
			return &types.QueryRateLimitsByChainIDResponse{}, errorsmod.Wrapf(types.ErrInvalidClientState, "Unable to fetch client state from channelID")
		}
		client, ok := clientState.(*ibctmtypes.ClientState)
		if !ok {
			return &types.QueryRateLimitsByChainIDResponse{}, errorsmod.Wrapf(types.ErrInvalidClientState, "Client state is not tendermint")
		}

		// If the chain ID matches, add the rate limit to the returned list
		if client.ChainId == req.ChainId {
			rateLimits = append(rateLimits, rateLimit)
		}
	}

	return &types.QueryRateLimitsByChainIDResponse{RateLimits: rateLimits}, nil
}

// RateLimitsByChannelID queries all rate limits for a given channel.
func (q queryServer) RateLimitsByChannelID(c context.Context, req *types.QueryRateLimitsByChannelIDRequest) (*types.QueryRateLimitsByChannelIDResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	rateLimits := []types.RateLimit{}
	for _, rateLimit := range q.GetAllRateLimits(ctx) {
		// If the channel ID matches, add the rate limit to the returned list
		if rateLimit.Path.ChannelID == req.ChannelID {
			rateLimits = append(rateLimits, rateLimit)
		}
	}

	return &types.QueryRateLimitsByChannelIDResponse{RateLimits: rateLimits}, nil
}

// AllWhitelistedAddresses queries all whitelisted addresses.
func (q queryServer) AllWhitelistedAddresses(c context.Context, _ *types.QueryAllWhitelistedAddressesRequest) (*types.QueryAllWhitelistedAddressesResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)
	whitelistedAddresses := q.GetAllWhitelistedAddressPairs(ctx)
	return &types.QueryAllWhitelistedAddressesResponse{AddressPairs: whitelistedAddresses}, nil
}
