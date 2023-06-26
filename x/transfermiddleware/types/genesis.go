package types

import (
	fmt "fmt"
)

// DefaultGenesisState returns a GenesisState with "transfer" as the default PortID.
func DefaultGenesisState() *GenesisState {
	return &GenesisState{}
}

// ValidateGenesis validates the provided genesis state to ensure the
// expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	err := validateTokenInfos(data.TokenInfos)
	if err != nil {
		return err
	}

	return nil
}

func validateTokenInfos(infos []ParachainIBCTokenInfo) error {
	infoMap := make(map[string]bool, len(infos))

	for i := 0; i < len(infos); i++ {
		info := infos[i]

		err := info.ValidateBasic()
		if err != nil {
			return err
		}

		strKey := info.AssetId

		// check duplicate based on assetId
		if _, ok := infoMap[strKey]; ok {
			return fmt.Errorf("duplicate parachain token info in genesis state: assetId %v, nativeDenom %v", info.AssetId, info.NativeDenom)
		}

		infoMap[strKey] = true
	}

	return nil
}
