package eventstream

import "time"

const (
	eventBufferSize   = 1024
	heartbeatInterval = 15 * time.Second
)
