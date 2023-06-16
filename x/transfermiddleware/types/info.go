package types

import (
	fmt "fmt"
	"strconv"
)

func (p ParachainIBCTokenInfo) ValidateBasic() error {
	_, err := strconv.Atoi(p.AssetId)
	if err != nil {
		return fmt.Errorf("Error parsing into int %v", p.AssetId)
	}

	return nil
}
