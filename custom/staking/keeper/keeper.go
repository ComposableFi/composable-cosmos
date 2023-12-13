package keeper

import (
	"github.com/cosmos/cosmos-sdk/codec"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	accountkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
	"github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

type Keeper struct {
	stakingkeeper.Keeper
	acck accountkeeper.AccountKeeper
}

var _ stakingkeeper.Keeper = stakingkeeper.Keeper{} //???

func NewBaseKeeper(
	cdc codec.BinaryCodec,
	key storetypes.StoreKey,
	ak types.AccountKeeper,
	acck accountkeeper.AccountKeeper,
	bk bankkeeper.Keeper,
	authority string,
) Keeper {
	keeper := Keeper{
		Keeper: *stakingkeeper.NewKeeper(cdc, key, ak, bk, authority),
		acck:   acck,
	}
	return keeper
}

func NewBaseKeeper2(
	staking stakingkeeper.Keeper,
	acck accountkeeper.AccountKeeper,
) Keeper {
	keeper := Keeper{
		Keeper: staking,
		acck:   acck,
	}
	return keeper
}

// func (k *Keeper) RegisterKeepers(akk banktypes.StakingKeeper) {
// 	k.acck = sk
// }
