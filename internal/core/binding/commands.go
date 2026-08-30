package binding

import (
	"context"

	"github.com/lucasdillmann/nginx-ignition/internal/core/common/validation"
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
