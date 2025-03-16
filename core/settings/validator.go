package settings

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/value_range"
	"dillmann.com.br/nginx-ignition/core/host"
	"strconv"
)

const (
	maximumDefaultContentTypeLength = 128
)

var (
	timeoutRange           = value_range.New(1, int(^uint(0)>>1))
	intervalRange          = value_range.New(1, int(^uint(0)>>1))
	logLinesRange          = value_range.New(0, 10000)
	workerProcessesRange   = value_range.New(1, 100)
	workerConnectionsRange = value_range.New(32, 4096)
	maximumBodySizeRange   = value_range.New(1, int(^uint(0)>>1))
)

type validator struct {
	validateBindingCommand *host.ValidateBindingCommand
	delegate               *validation.ConsistencyValidator
}

func newValidator(validateBindingCommand *host.ValidateBindingCommand) *validator {
	return &validator{
		validateBindingCommand,
		validation.NewValidator(),
	}
}

func (v *validator) validate(ctx context.Context, settings *Settings) error {
	v.validateNginx(settings.Nginx)
	v.validateLogRotation(settings.LogRotation)
	v.validateCertificateAutoRenew(settings.CertificateAutoRenew)

	if err := v.validateGlobalBindings(ctx, settings.GlobalBindings); err != nil {
		return err
	}

	return nil
}

func (v *validator) validateNginx(settings *NginxSettings) {
	v.checkRange(settings.Timeouts.Read, timeoutRange, "nginx.timeouts.read")
	v.checkRange(settings.Timeouts.Send, timeoutRange, "nginx.timeouts.send")
	v.checkRange(settings.Timeouts.Connect, timeoutRange, "nginx.timeouts.connect")
	v.checkRange(settings.Timeouts.Keepalive, timeoutRange, "nginx.timeouts.keepalive")
	v.checkRange(settings.WorkerProcesses, workerProcessesRange, "nginx.workerProcesses")
	v.checkRange(settings.WorkerConnections, workerConnectionsRange, "nginx.workerConnections")
	v.checkRange(settings.MaximumBodySizeMb, maximumBodySizeRange, "nginx.maximumBodySizeMb")

	if settings.DefaultContentType == "" {
		v.delegate.Add("nginx.defaultContentType", "A value is required")
	}

	if len(settings.DefaultContentType) > maximumDefaultContentTypeLength {
		v.delegate.Add("nginx.defaultContentType", "Cannot have more than 128 characters")
	}
}

func (v *validator) validateLogRotation(settings *LogRotationSettings) {
	v.checkRange(settings.IntervalUnitCount, intervalRange, "logRotation.intervalUnitCount")
	v.checkRange(settings.MaximumLines, logLinesRange, "logRotation.maximumLines")
}

func (v *validator) validateCertificateAutoRenew(settings *CertificateAutoRenewSettings) {
	v.checkRange(settings.IntervalUnitCount, intervalRange, "certificateAutoRenew.intervalUnitCount")
}

func (v *validator) validateGlobalBindings(ctx context.Context, settings []*host.Binding) error {
	for index, binding := range settings {
		if err := (*v.validateBindingCommand)(ctx, "globalBindings", index, binding, v.delegate); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) checkRange(value int, r *value_range.ValueRange, path string) {
	if !r.Contains(value) {
		v.delegate.Add(path, "Must be between "+strconv.Itoa(r.Min)+" and "+strconv.Itoa(r.Max))
	}
}
