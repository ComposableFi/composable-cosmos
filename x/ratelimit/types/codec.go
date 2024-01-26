package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/legacy"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"

	sdk "github.com/cosmos/cosmos-sdk/types"

	legacyMsg "github.com/notional-labs/composable/v6/x/ratelimit/types/legacy"
)

// RegisterLegacyAminoCodec registers the account interfaces and concrete types on the
// provided LegacyAmino codec. These types are used for Amino JSON serialization
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	legacy.RegisterAminoMsg(cdc, &MsgAddRateLimit{}, "composable/MsgAddRateLimit")
	legacy.RegisterAminoMsg(cdc, &MsgUpdateRateLimit{}, "composable/MsgUpdateRateLimit")
	legacy.RegisterAminoMsg(cdc, &MsgRemoveRateLimit{}, "composable/MsgRemoveRateLimit")
	legacy.RegisterAminoMsg(cdc, &MsgResetRateLimit{}, "composable/MsgResetRateLimit")
}

func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddRateLimit{},
		&MsgUpdateRateLimit{},
		&MsgRemoveRateLimit{},
		&MsgResetRateLimit{},

		&legacyMsg.MsgAddRateLimit{},
		&legacyMsg.MsgUpdateRateLimit{},
		&legacyMsg.MsgRemoveRateLimit{},
		&legacyMsg.MsgResetRateLimit{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &legacyMsg.Msg_serviceDesc)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino     = codec.NewLegacyAmino()
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	sdk.RegisterLegacyAminoCodec(amino)

	// Register all Amino interfaces and concrete types on the authz  and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}
