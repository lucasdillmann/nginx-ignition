package binding

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type Commands struct {
	Validate func(
		ctx context.Context,
		path string,
		index int,
		binding *Binding,
		validationCtx *validation.ConsistencyValidator,
	) error
}
