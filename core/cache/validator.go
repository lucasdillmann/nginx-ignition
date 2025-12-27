package cache

import (
	"fmt"
	"path/filepath"
	"strings"

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

func (v *validator) validate(c *Cache) error {
	v.validateBasicSettings(c)
	v.validateStoragePath(c.StoragePath)
	v.validateConcurrencyLock(c.ConcurrencyLock)
	v.validateCollections(c)
	v.validateFileExtensions(c.FileExtensions)

	return v.delegate.Result()
}

func (v *validator) validateBasicSettings(c *Cache) {
	if strings.TrimSpace(c.Name) == "" {
		v.delegate.Add("name", validation.ValueMissingMessage)
	}

	if c.InactiveSeconds != nil && *c.InactiveSeconds < 1 {
		v.delegate.Add("inactiveSeconds", validation.ValueCannotBeZeroMessage)
	}

	if c.MaximumSizeMB != nil && *c.MaximumSizeMB < 1 {
		v.delegate.Add("maximumSizeMb", validation.ValueCannotBeZeroMessage)
	}

	if c.MinimumUsesBeforeCaching < 1 {
		v.delegate.Add("minimumUsesBeforeCaching", validation.ValueCannotBeZeroMessage)
	}
}

func (v *validator) validateStoragePath(path *string) {
	if path == nil || strings.TrimSpace(*path) == "" {
		return
	}

	trimmedPath := strings.TrimSpace(*path)
	if !filepath.IsAbs(trimmedPath) {
		v.delegate.Add("storagePath", "Value must be an absolute path")
	}
}

func (v *validator) validateConcurrencyLock(lock ConcurrencyLock) {
	if !lock.Enabled {
		return
	}

	if lock.TimeoutSeconds == nil {
		v.delegate.Add("concurrencyLock.timeoutSeconds", validation.ValueMissingMessage)
	} else if *lock.TimeoutSeconds < 1 {
		v.delegate.Add("concurrencyLock.timeoutSeconds", validation.ValueCannotBeZeroMessage)
	}

	if lock.AgeSeconds == nil {
		v.delegate.Add("concurrencyLock.ageSeconds", validation.ValueMissingMessage)
	} else if *lock.AgeSeconds < 1 {
		v.delegate.Add("concurrencyLock.ageSeconds", validation.ValueCannotBeZeroMessage)
	}
}

func (v *validator) validateCollections(c *Cache) {
	for index, method := range c.AllowedMethods {
		v.validateMethod(index, method)
	}

	for index, staleOption := range c.UseStale {
		v.validateUseStaleOption(index, staleOption)
	}

	for index, duration := range c.Durations {
		v.validateDuration(index, &duration)
	}
}

func (v *validator) validateMethod(index int, method Method) {
	switch method {
	case GetMethod, HeadMethod, PostMethod, PutMethod, DeleteMethod, PatchMethod, OptionsMethod:
		// Valid
	default:
		path := fmt.Sprintf("allowedMethods[%d]", index)
		v.delegate.Add(path, "Invalid HTTP method")
	}
}

func (v *validator) validateUseStaleOption(index int, option UseStaleOption) {
	switch option {
	case ErrorUseStale, TimeoutUseStale, InvalidHeaderUseStale, UpdatingUseStale,
		Http500UseStale, Http502UseStale, Http503UseStale, Http504UseStale,
		Http403UseStale, Http404UseStale, Http429UseStale:
		// Valid
	default:
		path := fmt.Sprintf("useStale[%d]", index)
		v.delegate.Add(path, "Invalid stale cache option")
	}
}

func (v *validator) validateDuration(index int, duration *Duration) {
	path := fmt.Sprintf("durations[%d]", index)
	if len(duration.StatusCodes) == 0 {
		v.delegate.Add(path+".statusCodes", validation.ValueMissingMessage)
	}

	for statusCodeIndex, statusCode := range duration.StatusCodes {
		if httpStatusCodeRange.Contains(statusCode) {
			continue
		}

		v.delegate.Add(
			fmt.Sprintf("%s.statusCodes[%d]", path, statusCodeIndex),
			"Invalid status code (must be between 100 and 599 inclusive)",
		)
	}

	if duration.ValidTimeSeconds < 1 {
		v.delegate.Add(path+".validTimeSeconds", validation.ValueCannotBeZeroMessage)
	}
}

func (v *validator) validateFileExtensions(extensions []string) {
	for index, extension := range extensions {
		path := fmt.Sprintf("fileExtensions[%d]", index)

		if strings.TrimSpace(extension) == "" {
			v.delegate.Add(path, validation.ValueMissingMessage)
			continue
		}

		if strings.HasPrefix(extension, ".") {
			v.delegate.Add(path, "File extension cannot start with a dot")
		}
	}
}
