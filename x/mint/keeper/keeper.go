package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/notional-labs/composable/v6/x/mint/types"
)

// Keeper of the mint store
type Keeper struct {
	cdc              codec.BinaryCodec
	storeKey         storetypes.StoreKey
	accountKeeper    types.AccountKeeper
	stakingKeeper    types.StakingKeeper
	bankKeeper       types.BankKeeper
	feeCollectorName string

	// the address capable of executing a MsgUpdateParams message. Typically, this
	// should be the x/gov module account.
	authority string
}

// NewKeeper creates a new mint Keeper instance
func NewKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	sk types.StakingKeeper,
	ak types.AccountKeeper,
	bk types.BankKeeper,
	feeCollectorName string,
	authority string,
) Keeper {
	// ensure mint module account is set
	if addr := ak.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("the x/%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		cdc:              cdc,
		storeKey:         key,
		stakingKeeper:    sk,
		bankKeeper:       bk,
		feeCollectorName: feeCollectorName,
		authority:        authority,
		accountKeeper:    ak,
	}
}

func (k Keeper) GetModuleAccountAccAddress(ctx sdk.Context) sdk.AccAddress {
	moduleAccount := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	return moduleAccount.GetAddress()
}

// GetAuthority returns the x/mint module's authority.
func (k Keeper) GetAuthority() string {
	return k.authority
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", "x/"+types.ModuleName)
}

// GetMinter returns the minter.
func (k Keeper) GetMinter(ctx sdk.Context) (minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.MinterKey)
	if bz == nil {
		panic("stored minter should not have been nil")
	}

	if err := k.cdc.Unmarshal(bz, &minter); err != nil {
		panic(fmt.Sprintf("failed to unmarshal minter: %v", err))
	}
	return
}

// SetMinter sets the minter.
func (k Keeper) SetMinter(ctx sdk.Context, minter types.Minter) {
	store := ctx.KVStore(k.storeKey)
	bz := k.cdc.MustMarshal(&minter)
	store.Set(types.MinterKey, bz)
}

// SetAllowedAddress sets fund allowed address to state.
func (k Keeper) SetAllowedAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllowedAddressStoreKey(address)
	store.Set(key, []byte{1})
}

// IsAllowedAddress return true if has address in store
func (k Keeper) IsAllowedAddress(ctx sdk.Context, address string) bool {
	store := ctx.KVStore(k.storeKey)
	key := types.GetAllowedAddressStoreKey(address)
	return store.Has(key)
}

// SetParams sets the x/mint module parameters.
func (k Keeper) SetParams(ctx sdk.Context, p types.Params) error {
	if err := p.Validate(); err != nil {
		return fmt.Errorf("invalid params: %w", err)
	}

	store := ctx.KVStore(k.storeKey)
	bz, err := k.cdc.Marshal(&p)
	if err != nil {
		return fmt.Errorf("failed to marshal params: %w", err)
	}

	store.Set(types.ParamsKey, bz)
	return nil
}

// GetParams returns the current x/mint module parameters.
func (k Keeper) GetParams(ctx sdk.Context) (p types.Params) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.ParamsKey)
	if bz == nil {
		return p
	}

	k.cdc.MustUnmarshal(bz, &p)
	return p
}

// StakingTokenSupply implements an alias call to the underlying staking keeper's
// StakingTokenSupply to be used in BeginBlocker.
func (k Keeper) StakingTokenSupply(ctx sdk.Context) math.Int {
	return k.stakingKeeper.StakingTokenSupply(ctx)
}

// BondedRatio implements an alias call to the underlying staking keeper's
// BondedRatio to be used in BeginBlocker.
func (k Keeper) BondedRatio(ctx sdk.Context) math.LegacyDec {
	return k.stakingKeeper.BondedRatio(ctx)
}

// MintCoins implements an alias call to the underlying supply keeper's
// MintCoins to be used in BeginBlocker.
func (k Keeper) MintCoins(ctx sdk.Context, newCoins sdk.Coins) error {
	if newCoins.Empty() {
		// skip as no coins need to be minted
		return nil
	}

	return k.bankKeeper.MintCoins(ctx, types.ModuleName, newCoins)
}

// GetProvisionsFromBlock returns the provisions for a block based on the mint params
func (k Keeper) GetProvisionsFromBlock(ctx sdk.Context) (sdk.Coin, error) {
	minter := k.GetMinter(ctx)
	params := k.GetParams(ctx)

	// BlockProvisionWithCheck already includes the check for BlocksPerYear
	return k.BlockProvisionWithCheck(ctx, minter, params)
}

// BeginBlocker mints new tokens for the previous block.
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	provisions, err := k.GetProvisionsFromBlock(ctx)
	if err != nil {
		k.Logger(ctx).Error("failed to get provisions from block", "error", err)
		return
	}

	if err := k.MintCoins(ctx, sdk.NewCoins(provisions)); err != nil {
		k.Logger(ctx).Error("failed to mint coins", "error", err)
		return
	}

	// send the minted coins to the fee collector account
	if err := k.DistributeMintedCoin(ctx, provisions); err != nil {
		k.Logger(ctx).Error("failed to distribute minted coins", "error", err)
		return
	}
}

// DistributeMintedCoin implements the distribution of minted coins to the fee collector account
func (k Keeper) DistributeMintedCoin(ctx sdk.Context, mintedCoin sdk.Coin) error {
	if mintedCoin.IsZero() {
		return nil
	}

	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, sdk.NewCoins(mintedCoin))
}

// AddCollectedFees implements an alias call to the underlying supply keeper's
// AddCollectedFees to be used in BeginBlocker.
func (k Keeper) AddCollectedFees(ctx sdk.Context, fees sdk.Coins) error {
	return k.bankKeeper.SendCoinsFromModuleToModule(ctx, types.ModuleName, k.feeCollectorName, fees)
}
