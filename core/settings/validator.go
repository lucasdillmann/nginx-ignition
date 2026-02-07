package settings

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
)

const (
	maximumDefaultContentTypeLength = 128
	maximumRuntimeUserLength        = 32
	maximumDatabaseLocationLength   = 128
	defaultContentTypePath          = "nginx.defaultContentType"
)

var (
	timeoutRange           = valuerange.New(1, int(^uint(0)>>1))
	intervalRange          = valuerange.New(1, int(^uint(0)>>1))
	logLinesRange          = valuerange.New(0, 99_999)
	workerProcessesRange   = valuerange.New(1, 100)
	workerConnectionsRange = valuerange.New(32, 4096)
	maximumBodySizeRange   = valuerange.New(1, int(^uint(0)>>1))
	statsMaximumSizeRange  = valuerange.New(1, 512)
)

type validator struct {
	commands binding.Commands
	delegate *validation.ConsistencyValidator
}

func newValidator(commands binding.Commands) *validator {
	return &validator{
		commands,
		validation.NewValidator(),
	}
}

func (v *validator) validate(ctx context.Context, settings *Settings) error {
	v.validateNginx(ctx, settings.Nginx)
	v.validateLogRotation(ctx, settings.LogRotation)
	v.validateCertificateAutoRenew(ctx, settings.CertificateAutoRenew)

	if err := v.validateGlobalBindings(ctx, settings.GlobalBindings); err != nil {
		return err
	}

	return v.delegate.Result()
}

func (v *validator) validateNginx(ctx context.Context, settings *NginxSettings) {
	v.checkRange(ctx, settings.Timeouts.Read, timeoutRange, "nginx.timeouts.read")
	v.checkRange(ctx, settings.Timeouts.Send, timeoutRange, "nginx.timeouts.send")
	v.checkRange(ctx, settings.Timeouts.Connect, timeoutRange, "nginx.timeouts.connect")
	v.checkRange(ctx, settings.Timeouts.Keepalive, timeoutRange, "nginx.timeouts.keepalive")
	v.checkRange(ctx, settings.WorkerProcesses, workerProcessesRange, "nginx.workerProcesses")
	v.checkRange(ctx, settings.WorkerConnections, workerConnectionsRange, "nginx.workerConnections")
	v.checkRange(ctx, settings.MaximumBodySizeMb, maximumBodySizeRange, "nginx.maximumBodySizeMb")
	v.validateStats(ctx, settings.Stats)

	if settings.DefaultContentType == "" {
		v.delegate.Add(defaultContentTypePath, i18n.M(ctx, i18n.K.CommonValueMissing))
	}

	if len(settings.DefaultContentType) > maximumDefaultContentTypeLength {
		v.delegate.Add(
			defaultContentTypePath,
			i18n.M(ctx, i18n.K.CommonValueTooLong).V("max", maximumDefaultContentTypeLength),
		)
	}

	if strings.TrimSpace(settings.RuntimeUser) == "" {
		v.delegate.Add("nginx.runtimeUser", i18n.M(ctx, i18n.K.CommonValueMissing))
	}

	if len(settings.RuntimeUser) > maximumRuntimeUserLength {
		v.delegate.Add(
			"nginx.runtimeUser",
			i18n.M(ctx, i18n.K.CommonValueTooLong).V("max", maximumRuntimeUserLength),
		)
	}
}

func (v *validator) validateStats(ctx context.Context, settings *NginxStatsSettings) {
	v.checkRange(ctx, settings.MaximumSizeMB, statsMaximumSizeRange, "nginx.stats.maximumSizeMb")

	if settings.DatabaseLocation != nil && *settings.DatabaseLocation != "" {
		location := *settings.DatabaseLocation
		if !strings.HasSuffix(strings.ToLower(location), ".db") {
			v.delegate.Add(
				"nginx.stats.databaseLocation",
				i18n.M(ctx, i18n.K.CoreSettingsInvalidExtension).V("extension", ".db"),
			)
		}

		dir := filepath.Dir(location)
		if info, err := os.Stat(dir); err != nil || !info.IsDir() {
			v.delegate.Add(
				"nginx.stats.databaseLocation",
				i18n.M(ctx, i18n.K.CoreSettingsInvalidFolder),
			)
		}

		if len(location) > maximumDatabaseLocationLength {
			v.delegate.Add(
				"nginx.stats.databaseLocation",
				i18n.M(ctx, i18n.K.CommonValueTooLong).V("max", maximumDatabaseLocationLength),
			)
		}
	}
}

func (v *validator) validateLogRotation(ctx context.Context, settings *LogRotationSettings) {
	v.checkRange(ctx, settings.IntervalUnitCount, intervalRange, "logRotation.intervalUnitCount")
	v.checkRange(ctx, settings.MaximumLines, logLinesRange, "logRotation.maximumLines")
}

func (v *validator) validateCertificateAutoRenew(
	ctx context.Context,
	settings *CertificateAutoRenewSettings,
) {
	v.checkRange(
		ctx,
		settings.IntervalUnitCount,
		intervalRange,
		"certificateAutoRenew.intervalUnitCount",
	)
}

func (v *validator) validateGlobalBindings(ctx context.Context, settings []binding.Binding) error {
	for index, b := range settings {
		if err := v.commands.Validate(ctx, "globalBindings", index, &b, v.delegate); err != nil {
			return err
		}
	}

	return nil
}

func (v *validator) checkRange(
	ctx context.Context,
	value int,
	r *valuerange.ValueRange,
	path string,
) {
	if !r.Contains(value) {
		v.delegate.Add(
			path,
			i18n.M(ctx, i18n.K.CommonBetweenValues).V("min", r.Min).V("max", r.Max),
		)
	}
}
