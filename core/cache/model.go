package cache

import (
	"github.com/google/uuid"
)

type Method string

const (
	GetMethod     Method = "GET"
	HeadMethod    Method = "HEAD"
	PostMethod    Method = "POST"
	PutMethod     Method = "PUT"
	DeleteMethod  Method = "DELETE"
	PatchMethod   Method = "PATCH"
	OptionsMethod Method = "OPTIONS"
)

type UseStaleOption string

const (
	ErrorUseStale         UseStaleOption = "ERROR"
	TimeoutUseStale       UseStaleOption = "TIMEOUT"
	InvalidHeaderUseStale UseStaleOption = "INVALID_HEADER"
	UpdatingUseStale      UseStaleOption = "UPDATING"
	HTTP500UseStale       UseStaleOption = "HTTP_500"
	HTTP502UseStale       UseStaleOption = "HTTP_502"
	HTTP503UseStale       UseStaleOption = "HTTP_503"
	HTTP504UseStale       UseStaleOption = "HTTP_504"
	HTTP403UseStale       UseStaleOption = "HTTP_403"
	HTTP404UseStale       UseStaleOption = "HTTP_404"
	HTTP429UseStale       UseStaleOption = "HTTP_429"
)

type Cache struct {
	InactiveSeconds                  *int
	StoragePath                      *string
	MaximumSizeMB                    *int
	ConcurrencyLock                  ConcurrencyLock
	Name                             string
	UseStale                         []UseStaleOption
	AllowedMethods                   []Method
	BypassRules                      []string
	NoCacheRules                     []string
	FileExtensions                   []string
	Durations                        []Duration
	MinimumUsesBeforeCaching         int
	ID                               uuid.UUID
	Revalidate                       bool
	BackgroundUpdate                 bool
	IgnoreUpstreamCacheHeaders       bool
	CacheStatusResponseHeaderEnabled bool
}

type ConcurrencyLock struct {
	TimeoutSeconds *int
	AgeSeconds     *int
	Enabled        bool
}

type Duration struct {
	StatusCodes      []string
	ValidTimeSeconds int
}
