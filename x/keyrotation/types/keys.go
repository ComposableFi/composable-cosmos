package types

const (
	ModuleName = "keyrotation"

	// QuerierRoute defines the module's query routing key
	QuerierRoute = ModuleName

	RouterKey = ModuleName
)

func KeyPrefix(p string) []byte {
	return []byte(p)
}

var (
	KeyRotation = KeyPrefix("rotation")
)

func GetKeyRotationHistory(valAddress string) []byte {
	return append(KeyRotation, []byte(valAddress)...)
}
