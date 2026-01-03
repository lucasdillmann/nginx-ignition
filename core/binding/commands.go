package binding

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type Commands interface {
	Validate(
		ctx context.Context,
		path string,
		index int,
		binding *Binding,
		validationCtx *validation.ConsistencyValidator,
	) error
}
