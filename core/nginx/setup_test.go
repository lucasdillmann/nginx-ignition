package nginx

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/container"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

func init() {
	_ = log.Init()
	container.Init(context.Background())
}
