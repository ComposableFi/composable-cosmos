package prepare

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/authz"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
)

// FilterChangeValSetMsgs slpit list txs to 2 set : filteredTxs, changeValSetTxs
func FilterChangeValSetMsgs(ctx sdk.Context, cdc codec.BinaryCodec, decoder sdk.TxDecoder, txs [][]byte) [][]byte {
	var filteredTxs [][]byte
	var changeValSetTxs [][]byte
	for _, txBytes := range txs {
		// Decode tx so we can read msgs.
		tx, err := decoder(txBytes)
		if err != nil {
			ctx.Logger().Error(fmt.Sprintf("RemoveDisallowMsgs: failed to decode tx: %v", err))
			continue // continue to next tx.
		}

		// For each msg in tx, check if it is disallowed.
		containsChangeValsetMsg := false
		for _, msg := range tx.GetMsgs() {
			if isChangeValSetMsg(ctx, cdc, msg) {
				containsChangeValsetMsg = true
				break // break out of loop over msgs.
			}
		}

		// If tx contains disallowed msg, skip it.
		if containsChangeValsetMsg {
			ctx.Logger().Info("Found change valset msg")
			changeValSetTxs = append(changeValSetTxs, txBytes)
			continue // continue to next tx.
		}

		// Otherwise, add tx to filtered txs.
		filteredTxs = append(filteredTxs, txBytes)
	}

	return filteredTxs
}

func isChangeValSetMsg(ctx sdk.Context, cdc codec.BinaryCodec, msg sdk.Msg) bool {
	switch msg := msg.(type) {

	case *stakingtypes.MsgDelegate:
		return true
	case *stakingtypes.MsgBeginRedelegate:
		return true
	case *authz.MsgExec:
		return isChangeValSetAuthzMsg(ctx, cdc, msg)
	default:
		return false
	}
}

func isChangeValSetAuthzMsg(ctx sdk.Context, cdc codec.BinaryCodec, execMsg *authz.MsgExec) bool {
	for _, v := range execMsg.Msgs {
		var innerMsg sdk.Msg
		if err := cdc.UnpackAny(v, &innerMsg); err != nil {
			return false
		}
		if isChangeValSetMsg(ctx, cdc, innerMsg) {
			return true
		}
	}
	return false
}
