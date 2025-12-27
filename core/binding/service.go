package binding

import (
	"context"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

type service struct {
	certificateCommands *certificate.Commands
}

func newService(certificateCommands *certificate.Commands) *service {
	return &service{certificateCommands}
}

func (s *service) validateBinding(
	ctx context.Context,
	path string,
	index int,
	binding *Binding,
	validationCtx *validation.ConsistencyValidator,
) error {
	return newValidator(validationCtx, s.certificateCommands).validate(ctx, path, binding, index)
}
