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
	ID                       uuid.UUID
	Name                     string
	StoragePath              *string
	InactiveSeconds          *int
	MaxSizeMB                *int
	AllowedMethods           []Method
	MinimumUsesBeforeCaching *int
	UseStale                 []UseStaleOption
	BackgroundUpdate         *bool
	ConcurrencyLock          ConcurrencyLock
	Revalidate               *bool
	BypassRules              []string
	NoCacheRules             []string
	Durations                []Duration
}

type ConcurrencyLock struct {
	Enabled        bool
	TimeoutSeconds *int
	AgeSeconds     *int
}

type Duration struct {
	StatusCodes      []int
	ValidTimeSeconds int
}
