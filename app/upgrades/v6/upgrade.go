package v6

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/module"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	"github.com/notional-labs/composable/v6/app/keepers"
	"github.com/notional-labs/composable/v6/app/upgrades"

	bech32authmigration "github.com/notional-labs/composable/v6/bech32-migration/auth"
	bech32govmigration "github.com/notional-labs/composable/v6/bech32-migration/gov"
	bech32slashingmigration "github.com/notional-labs/composable/v6/bech32-migration/slashing"
	bech32stakingmigration "github.com/notional-labs/composable/v6/bech32-migration/staking"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	"github.com/notional-labs/composable/v6/x/ratelimit/types"

	"github.com/notional-labs/composable/v6/bech32-migration/utils"
)

const (
	// https://github.com/cosmos/chain-registry/blob/master/composable/assetlist.json
	uatom = "ibc/EF48E6B1A1A19F47ECAEA62F5670C37C0580E86A9E88498B7E393EB6F49F33C0"
	dot   = "ibc/3CC19CEC7E5A3E90E78A5A9ECC5A0E2F8F826A375CF1E096F4515CF09DA3E366"
	ksm   = "ibc/EE9046745AEC0E8302CB7ED9D5AD67F528FB3B7AE044B247FB0FB293DBDA35E9"
	usdt  = "ibc/F3EC9F834E57DF704FA3AEAF14E8391C2E58397FE56960AD70E67562990D8265"
)

func CreateUpgradeHandler(
	mm *module.Manager,
	configurator module.Configurator,
	_ upgrades.BaseAppParamManager,
	cdc codec.Codec,
	keepers *keepers.AppKeepers,
) upgradetypes.UpgradeHandler {
	return func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {
		// Migration prefix
		ctx.Logger().Info("First step: Migrate addresses stored in bech32 form to use new prefix")
		keys := keepers.GetKVStoreKey()
		bech32stakingmigration.MigrateAddressBech32(ctx, keys[stakingtypes.StoreKey], cdc)
		bech32slashingmigration.MigrateAddressBech32(ctx, keys[slashingtypes.StoreKey], cdc)
		bech32govmigration.MigrateAddressBech32(ctx, keys[govtypes.StoreKey], cdc)
		bech32authmigration.MigrateAddressBech32(ctx, keys[authtypes.StoreKey], cdc)

		// allowed relayer address
		tfmdk := keepers.TransferMiddlewareKeeper
		tfmdk.IterateAllowRlyAddress(ctx, func(rlyAddress string) bool {
			// Delete old address
			tfmdk.DeleteAllowRlyAddress(ctx, rlyAddress)

			// add new address
			newRlyAddress := utils.ConvertAccAddr(rlyAddress)
			tfmdk.SetAllowRlyAddress(ctx, newRlyAddress)
			return false
		})

		// update rate limit to 50k
		rlKeeper := keepers.RatelimitKeeper
		// add uatom
		uatomRateLimit, found := rlKeeper.GetRateLimit(ctx, uatom, "channel-2")
		if !found {
			channelValue := rlKeeper.GetChannelValue(ctx, uatom)
			// Create and store the rate limit object
			path := types.Path{
				Denom:     uatom,
				ChannelID: "channel-2",
			}
			quota := types.Quota{
				MaxPercentSend: sdk.NewInt(30),
				MaxPercentRecv: sdk.NewInt(30),
				DurationHours:  24,
			}
			flow := types.Flow{
				Inflow:       math.ZeroInt(),
				Outflow:      math.ZeroInt(),
				ChannelValue: channelValue,
			}
			uatomRateLimit = types.RateLimit{
				Path:               &path,
				Quota:              &quota,
				Flow:               &flow,
				MinRateLimitAmount: sdk.NewInt(1282_000_000 * 5), // 1282_000_000 * 5
			}
			rlKeeper.SetRateLimit(ctx, uatomRateLimit)
		} else {
			uatomRateLimit.MinRateLimitAmount = sdk.NewInt(1282_000_000 * 5)
			rlKeeper.SetRateLimit(ctx, uatomRateLimit)
		}
		// add dot
		dotRateLimit, found := rlKeeper.GetRateLimit(ctx, dot, "channel-2")
		if !found {
			channelValue := rlKeeper.GetChannelValue(ctx, dot)
			// Create and store the rate limit object
			path := types.Path{
				Denom:     dot,
				ChannelID: "channel-2",
			}
			quota := types.Quota{
				MaxPercentSend: sdk.NewInt(30),
				MaxPercentRecv: sdk.NewInt(30),
				DurationHours:  24,
			}
			flow := types.Flow{
				Inflow:       math.ZeroInt(),
				Outflow:      math.ZeroInt(),
				ChannelValue: channelValue,
			}
			dotRateLimit = types.RateLimit{
				Path:               &path,
				Quota:              &quota,
				Flow:               &flow,
				MinRateLimitAmount: sdk.NewInt(22_670_000_000_000 * 5), // 22_670_000_000_000 * 5
			}
			rlKeeper.SetRateLimit(ctx, dotRateLimit)
		} else {
			dotRateLimit.MinRateLimitAmount = sdk.NewInt(22_670_000_000_000 * 5)
			rlKeeper.SetRateLimit(ctx, dotRateLimit)
		}
		// add ksm
		ksmRateLimit, found := rlKeeper.GetRateLimit(ctx, ksm, "channel-2")
		if !found {
			channelValue := rlKeeper.GetChannelValue(ctx, ksm)
			// Create and store the rate limit object
			path := types.Path{
				Denom:     ksm,
				ChannelID: "channel-2",
			}
			quota := types.Quota{
				MaxPercentSend: sdk.NewInt(30),
				MaxPercentRecv: sdk.NewInt(30),
				DurationHours:  24,
			}
			flow := types.Flow{
				Inflow:       math.ZeroInt(),
				Outflow:      math.ZeroInt(),
				ChannelValue: channelValue,
			}
			ksmRateLimit = types.RateLimit{
				Path:               &path,
				Quota:              &quota,
				Flow:               &flow,
				MinRateLimitAmount: sdk.NewInt(510_000_000_000_000 * 5), // 510_000_000_000_000*5
			}
			rlKeeper.SetRateLimit(ctx, ksmRateLimit)
		} else {
			ksmRateLimit.MinRateLimitAmount = sdk.NewInt(510_000_000_000_000 * 5)
			rlKeeper.SetRateLimit(ctx, ksmRateLimit)
		}
		// add usdt
		usdtRateLimit, found := rlKeeper.GetRateLimit(ctx, usdt, "channel-2")
		if !found {
			channelValue := rlKeeper.GetChannelValue(ctx, usdt)
			// Create and store the rate limit object
			path := types.Path{
				Denom:     usdt,
				ChannelID: "channel-2",
			}
			quota := types.Quota{
				MaxPercentSend: sdk.NewInt(30),
				MaxPercentRecv: sdk.NewInt(30),
				DurationHours:  24,
			}
			flow := types.Flow{
				Inflow:       math.ZeroInt(),
				Outflow:      math.ZeroInt(),
				ChannelValue: channelValue,
			}
			usdtRateLimit = types.RateLimit{
				Path:               &path,
				Quota:              &quota,
				Flow:               &flow,
				MinRateLimitAmount: sdk.NewInt(10_000_000_000 * 5), // 10_000_000_000 * 5
			}
			rlKeeper.SetRateLimit(ctx, usdtRateLimit)
		} else {
			usdtRateLimit.MinRateLimitAmount = sdk.NewInt(10_000_000_000 * 5)
			rlKeeper.SetRateLimit(ctx, usdtRateLimit)
		}

		return mm.RunMigrations(ctx, configurator, vm)
	}
}
