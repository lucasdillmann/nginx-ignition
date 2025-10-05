package dns

import (
	"time"
)

const (
	PropagationTimeout = 180 * time.Second
	SequenceInterval   = 180 * time.Second
	PollingInterval    = 1 * time.Second
	TTL                = 300
	MaxRetries         = 3
)
