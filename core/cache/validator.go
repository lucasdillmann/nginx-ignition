package cache

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/common/valuerange"
)

var httpStatusCodeRange = valuerange.New(100, 599)

type validator struct {
	delegate *validation.ConsistencyValidator
}

func newValidator() *validator {
	return &validator{
		delegate: validation.NewValidator(),
	}
}

func (v *validator) validate(ctx context.Context, c *Cache) error {
	v.validateBasicSettings(ctx, c)
	v.validateStoragePath(ctx, c.StoragePath)
	v.validateConcurrencyLock(ctx, c.ConcurrencyLock)
	v.validateCollections(ctx, c)
	v.validateFileExtensions(ctx, c.FileExtensions)

	return v.delegate.Result()
}

func (v *validator) validateBasicSettings(ctx context.Context, c *Cache) {
	if strings.TrimSpace(c.Name) == "" {
		v.delegate.Add("name", i18n.M(ctx, "common.validation.value-missing"))
	}

	if c.InactiveSeconds != nil && *c.InactiveSeconds < 1 {
		v.delegate.Add("inactiveSeconds", i18n.M(ctx, "common.validation.cannot-be-zero"))
	}

	if c.MaximumSizeMB != nil && *c.MaximumSizeMB < 1 {
		v.delegate.Add("maximumSizeMb", i18n.M(ctx, "common.validation.cannot-be-zero"))
	}

	if c.MinimumUsesBeforeCaching < 1 {
		v.delegate.Add("minimumUsesBeforeCaching", i18n.M(ctx, "common.validation.cannot-be-zero"))
	}
}

func (v *validator) validateStoragePath(ctx context.Context, path *string) {
	if path == nil || strings.TrimSpace(*path) == "" {
		return
	}

	trimmedPath := strings.TrimSpace(*path)
	if !filepath.IsAbs(trimmedPath) {
		v.delegate.Add("storagePath", i18n.M(ctx, "common.validation.absolute-path-required"))
	}
}

func (v *validator) validateConcurrencyLock(ctx context.Context, lock ConcurrencyLock) {
	if !lock.Enabled {
		return
	}

	if lock.TimeoutSeconds == nil {
		v.delegate.Add(
			"concurrencyLock.timeoutSeconds",
			i18n.M(ctx, "common.validation.value-missing"),
		)
	} else if *lock.TimeoutSeconds < 1 {
		v.delegate.Add(
			"concurrencyLock.timeoutSeconds",
			i18n.M(ctx, "common.validation.cannot-be-zero"),
		)
	}

	if lock.AgeSeconds == nil {
		v.delegate.Add("concurrencyLock.ageSeconds", i18n.M(ctx, "common.validation.value-missing"))
	} else if *lock.AgeSeconds < 1 {
		v.delegate.Add(
			"concurrencyLock.ageSeconds",
			i18n.M(ctx, "common.validation.cannot-be-zero"),
		)
	}
}

func (v *validator) validateCollections(ctx context.Context, c *Cache) {
	for index, method := range c.AllowedMethods {
		v.validateMethod(ctx, index, method)
	}

	for index, staleOption := range c.UseStale {
		v.validateUseStaleOption(ctx, index, staleOption)
	}

	for index, duration := range c.Durations {
		v.validateDuration(ctx, index, &duration)
	}
}

func (v *validator) validateMethod(ctx context.Context, index int, method Method) {
	switch method {
	case GetMethod, HeadMethod, PostMethod, PutMethod, DeleteMethod, PatchMethod, OptionsMethod:
		// Valid
	default:
		path := fmt.Sprintf("allowedMethods[%d]", index)
		v.delegate.Add(path, i18n.M(ctx, "cache.validation.invalid-method"))
	}
}

func (v *validator) validateUseStaleOption(ctx context.Context, index int, option UseStaleOption) {
	switch option {
	case ErrorUseStale, TimeoutUseStale, InvalidHeaderUseStale, UpdatingUseStale,
		HTTP500UseStale, HTTP502UseStale, HTTP503UseStale, HTTP504UseStale,
		HTTP403UseStale, HTTP404UseStale, HTTP429UseStale:
		// Valid
	default:
		path := fmt.Sprintf("useStale[%d]", index)
		v.delegate.Add(path, i18n.M(ctx, "cache.validation.invalid-stale-option"))
	}
}

func (v *validator) validateDuration(ctx context.Context, index int, duration *Duration) {
	path := fmt.Sprintf("durations[%d].statusCodes", index)
	if len(duration.StatusCodes) == 0 {
		v.delegate.Add(path, i18n.M(ctx, "common.validation.value-missing"))
	}

	for _, rawValue := range duration.StatusCodes {
		statusCode, err := strconv.Atoi(rawValue)
		if err == nil && httpStatusCodeRange.Contains(statusCode) {
			continue
		}

		v.delegate.Add(
			path,
			i18n.M(ctx, "cache.validation.invalid-status-code").
				V("value", rawValue).
				V("min", httpStatusCodeRange.Min).
				V("max", httpStatusCodeRange.Max),
		)
	}

	if duration.ValidTimeSeconds < 1 {
		v.delegate.Add(path+".validTimeSeconds", i18n.M(ctx, "common.validation.cannot-be-zero"))
	}
}

func (v *validator) validateFileExtensions(ctx context.Context, extensions []string) {
	for index, extension := range extensions {
		path := fmt.Sprintf("fileExtensions[%d]", index)

		if strings.TrimSpace(extension) == "" {
			v.delegate.Add(path, i18n.M(ctx, "common.validation.value-missing"))
			continue
		}

		if strings.HasPrefix(extension, ".") {
			v.delegate.Add(path, i18n.M(ctx, "cache.validation.extension-dot-not-allowed"))
		}
	}
}
