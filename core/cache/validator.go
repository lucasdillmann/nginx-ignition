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
	if strings.TrimSpace(c.Name) == "" {
		v.delegate.Add("name", validation.ValueMissingMessage)
	}

	if c.StoragePath != nil && strings.TrimSpace(*c.StoragePath) != "" {
		trimmedPath := strings.TrimSpace(*c.StoragePath)
		if !filepath.IsAbs(trimmedPath) {
			v.delegate.Add("storagePath", "Value must be an absolute path")
		} else if _, err := os.Stat(trimmedPath); os.IsNotExist(err) {
			v.delegate.Add("storagePath", "Path does not exist")
		}
	}

	if c.InactiveSeconds != nil && *c.InactiveSeconds < 0 {
		v.delegate.Add("inactiveSeconds", "Value must be 0 or greater")
	}

	if c.MaxSizeMB != nil && *c.MaxSizeMB < 1 {
		v.delegate.Add("maxSizeMB", "Value must be 1 or greater")
	}

	for index, method := range c.AllowedMethods {
		v.validateMethod(index, method)
	}

	if c.MinimumUsesBeforeCaching != nil && *c.MinimumUsesBeforeCaching < 1 {
		v.delegate.Add("minimumUsesBeforeCaching", "Value must be 1 or greater")
	}

	for index, staleOption := range c.UseStale {
		v.validateUseStaleOption(index, staleOption)
	}

	if c.ConcurrencyLock.Enabled {
		if c.ConcurrencyLock.TimeoutSeconds == nil {
			v.delegate.Add("concurrencyLock.timeoutSeconds", validation.ValueMissingMessage)
		} else if *c.ConcurrencyLock.TimeoutSeconds < 0 {
			v.delegate.Add("concurrencyLock.timeoutSeconds", "Value must be 0 or greater")
		}

		if c.ConcurrencyLock.AgeSeconds == nil {
			v.delegate.Add("concurrencyLock.ageSeconds", validation.ValueMissingMessage)
		} else if *c.ConcurrencyLock.AgeSeconds < 0 {
			v.delegate.Add("concurrencyLock.ageSeconds", "Value must be 0 or greater")
		}
	}

	for index, duration := range c.Durations {
		v.validateDuration(index, &duration)
	}

	return v.delegate.Result()
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
		v.delegate.Add(path+".validTimeSeconds", "Value must be 0 or greater")
	}
}
