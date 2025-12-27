package cache

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/validation"
)

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

	return v.delegate.Result()
}

func (v *validator) validateBasicSettings(c *Cache) {
	if strings.TrimSpace(c.Name) == "" {
		v.delegate.Add("name", validation.ValueMissingMessage)
	}

	if c.InactiveSeconds != nil && *c.InactiveSeconds < 0 {
		v.delegate.Add("inactiveSeconds", validation.ValueCannotBeNegativeMessage)
	}

	if c.MaximumSizeMB != nil && *c.MaximumSizeMB < 1 {
		v.delegate.Add("maximumSizeMb", "Value must be 1 or greater")
	}

	if c.MinimumUsesBeforeCaching < 1 {
		v.delegate.Add("minimumUsesBeforeCaching", "Value must be 1 or greater")
	}
}

func (v *validator) validateStoragePath(path *string) {
	if path == nil || strings.TrimSpace(*path) == "" {
		return
	}

	trimmedPath := strings.TrimSpace(*path)
	if !filepath.IsAbs(trimmedPath) {
		v.delegate.Add("storagePath", "Value must be an absolute path")
	} else if _, err := os.Stat(trimmedPath); os.IsNotExist(err) {
		v.delegate.Add("storagePath", "Path does not exist")
	}
}

func (v *validator) validateConcurrencyLock(lock ConcurrencyLock) {
	if !lock.Enabled {
		return
	}

	if lock.TimeoutSeconds == nil {
		v.delegate.Add("concurrencyLock.timeoutSeconds", validation.ValueMissingMessage)
	} else if *lock.TimeoutSeconds < 0 {
		v.delegate.Add("concurrencyLock.timeoutSeconds", validation.ValueCannotBeNegativeMessage)
	}

	if lock.AgeSeconds == nil {
		v.delegate.Add("concurrencyLock.ageSeconds", validation.ValueMissingMessage)
	} else if *lock.AgeSeconds < 0 {
		v.delegate.Add("concurrencyLock.ageSeconds", validation.ValueCannotBeNegativeMessage)
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
		path := "allowedMethods[" + strconv.Itoa(index) + "]"
		v.delegate.Add(path, "Invalid HTTP method")
	}
}

func (v *validator) validateUseStaleOption(index int, option UseStaleOption) {
	switch option {
	case ErrorUseStale, TimeoutUseStale, InvalidHeaderUseStale, UpdatingUseStale,
		Http500UseStale, Http502UseStale, Http503UseStale, Http504UseStale,
		Http403UseStale, Http404UseStale, Http429UseStale, OffUseStale:
		// Valid
	default:
		path := "useStale[" + strconv.Itoa(index) + "]"
		v.delegate.Add(path, "Invalid stale cache option")
	}
}

func (v *validator) validateDuration(index int, duration *Duration) {
	path := "durations[" + strconv.Itoa(index) + "]"
	if len(duration.StatusCodes) == 0 {
		v.delegate.Add(path+".statusCodes", validation.ValueMissingMessage)
	}

	if duration.ValidTimeSeconds < 0 {
		v.delegate.Add(path+".validTimeSeconds", validation.ValueCannotBeNegativeMessage)
	}
}
