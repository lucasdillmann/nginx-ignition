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
	ErrorUseStale         UseStaleOption = "error"
	TimeoutUseStale       UseStaleOption = "timeout"
	InvalidHeaderUseStale UseStaleOption = "invalid_header"
	UpdatingUseStale      UseStaleOption = "updating"
	Http500UseStale       UseStaleOption = "http_500"
	Http502UseStale       UseStaleOption = "http_502"
	Http503UseStale       UseStaleOption = "http_503"
	Http504UseStale       UseStaleOption = "http_504"
	Http403UseStale       UseStaleOption = "http_403"
	Http404UseStale       UseStaleOption = "http_404"
	Http429UseStale       UseStaleOption = "http_429"
	OffUseStale           UseStaleOption = "off"
)

type Cache struct {
	ConcurrencyLock          ConcurrencyLock
	InactiveSeconds          *int
	StoragePath              *string
	MaxSizeMB                *int
	MinimumUsesBeforeCaching *int
	BackgroundUpdate         *bool
	Revalidate               *bool
	Name                     string
	AllowedMethods           []Method
	UseStale                 []UseStaleOption
	BypassRules              []string
	NoCacheRules             []string
	Durations                []Duration
	ID                       uuid.UUID
}

type ConcurrencyLock struct {
	TimeoutSeconds *int
	AgeSeconds     *int
	Enabled        bool
}

type Duration struct {
	StatusCodes      []int
	ValidTimeSeconds int
}
