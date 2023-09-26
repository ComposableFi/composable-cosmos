package v5_2_0

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cosmos/ibc-go/v7/modules/core/exported"
	ibckeeper "github.com/cosmos/ibc-go/v7/modules/core/keeper"
	wasm08keeper "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/keeper"
	wasm08types "github.com/cosmos/ibc-go/v7/modules/light-clients/08-wasm/types"

	"github.com/notional-labs/centauri/v5/app/keepers"
)

const (
	newWasmCodeID      = ""
	clientId           = "08-wasm-05"
	substituteClientId = "08-wasm-06"
)

func RunForkLogic(ctx sdk.Context, keepers *keepers.AppKeepers) {
	ctx.Logger().Info("Applying v5_2_0 upgrade" +
		"Upgrade 08-wasm contract",
	)

	UpdateWasmContract(ctx, keepers.IBCKeeper, keepers.Wasm08Keeper)
	ClientUpdate(ctx, keepers.IBCKeeper.Codec(), keepers.IBCKeeper, clientId, substituteClientId)
}

func UpdateWasmContract(ctx sdk.Context, ibckeeper *ibckeeper.Keeper, wasmKeeper wasm08keeper.Keeper) {
	unknownClientState, found := ibckeeper.ClientKeeper.GetClientState(ctx, clientId)
	if !found {
		panic("cannot update client with ID")
	}

	clientState, ok := unknownClientState.(*wasm08types.ClientState)
	if !ok {
		panic("cannot update client with ID")
	}

	clientState.CodeId = []byte(newWasmCodeID)

	ibckeeper.ClientKeeper.SetClientState(ctx, clientId, clientState)
}

func ClientUpdate(ctx sdk.Context, codec codec.BinaryCodec, ibckeeper *ibckeeper.Keeper, subjectClientId string, substituteClientId string) error {
	subjectClientState, found := ibckeeper.ClientKeeper.GetClientState(ctx, subjectClientId)
	if !found {
		panic("cannot update client with ID")
	}

	subjectClientStore := ibckeeper.ClientKeeper.ClientStore(ctx, subjectClientId)

	if status := ibckeeper.ClientKeeper.GetClientStatus(ctx, subjectClientState, subjectClientId); status == exported.Active {
		panic("cannot update client with ID")
	}

	substituteClientState, found := ibckeeper.ClientKeeper.GetClientState(ctx, substituteClientId)
	if !found {
		panic("cannot update client with ID")
	}

	if subjectClientState.GetLatestHeight().GTE(substituteClientState.GetLatestHeight()) {
		panic("cannot update client with ID")
	}

	substituteClientStore := ibckeeper.ClientKeeper.ClientStore(ctx, substituteClientId)

	if status := ibckeeper.ClientKeeper.GetClientStatus(ctx, substituteClientState, substituteClientId); status != exported.Active {
		panic("cannot update client with ID")
	}

	if err := subjectClientState.CheckSubstituteAndUpdateState(ctx, codec, subjectClientStore, substituteClientStore, substituteClientState); err != nil {
		panic("cannot update client with ID")
	}

	ctx.Logger().Info("client updated after hark fork passed", "client-id", subjectClientId)

	return nil
}
