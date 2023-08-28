package types

const (
	// ModuleName defines the module name
	ModuleName = "txboundary"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName
)

var (
	// DelegateBoundaryKey is the key to use for the keeper store.
	DelegateBoundaryKey = []byte{0x00}

	// RedelegateBoundaryKey is the key to use for the keeper store.
	RedelegateBoundaryKey = []byte{0x01}
)
