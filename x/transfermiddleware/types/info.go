package types

import (
	fmt "fmt"
	"strconv"

	pcktfrwdtypes "github.com/cosmos/ibc-apps/middleware/packet-forward-middleware/v7/packetforward/types"
)

var (
	_ pcktfrwdtypes.ParaChainIBCTokenInfo = ParachainIBCTokenInfo{}
)

func (p ParachainIBCTokenInfo) ValidateBasic() error {
	_, err := strconv.Atoi(p.AssetId)
	if err != nil {
		return fmt.Errorf("error parsing into int %v", p.AssetId)
	}

	return nil
}

// GetNativeDenom implements interface.
func (p ParachainIBCTokenInfo) GetNativeDenom() string {
	return p.NativeDenom
}

// GetIbcDenom implements interface.
func (p ParachainIBCTokenInfo) GetIbcDenom() string {
	return p.IbcDenom
}

// GetChannelID implements interface.
func (p ParachainIBCTokenInfo) GetChannelID() string {
	return p.ChannelID
}
