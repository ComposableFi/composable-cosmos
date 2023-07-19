package types

var (
	EventTransferDenied = "transfer_denied"

	EventRateLimitExceeded = "rate_limit_exceeded"

	AttributeKeyReason  = "reason"
	AttributeKeyModule  = "module"
	AttributeKeyAction  = "action"
	AttributeKeyDenom   = "denom"
	AttributeKeyChannel = "channel"
	AttributeKeyAmount  = "amount"
	AttributeKeyError   = "error"

	EventTypeEpochEnd       = "epoch_end" // TODO: need to clean up (not use)
	EventTypeEpochStart     = "epoch_start"
	AttributeEpochNumber    = "epoch_number"
	AttributeEpochStartTime = "start_time"
)
