package nginx

import (
	"context"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"path/filepath"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type statsResponse struct {
	ServerZones   map[string]statsZoneData            `json:"serverZones"`
	FilterZones   map[string]map[string]statsZoneData `json:"filterZones"`
	UpstreamZones map[string][]statsUpstreamZoneData  `json:"upstreamZones"`
	HostName      string                              `json:"hostName"`
	Connections   statsConnections                    `json:"connections"`
}

type statsConnections struct {
	Active   uint64 `json:"active"`
	Reading  uint64 `json:"reading"`
	Writing  uint64 `json:"writing"`
	Waiting  uint64 `json:"waiting"`
	Accepted uint64 `json:"accepted"`
	Handled  uint64 `json:"handled"`
	Requests uint64 `json:"requests"`
}

type statsZoneData struct {
	RequestMsecs       statsTimeSeries `json:"requestMsecs"`
	RequestBuckets     statsBuckets    `json:"requestBuckets"`
	OverCounts         statsOverCounts `json:"overCounts"`
	Responses          statsResponses  `json:"responses"`
	RequestCounter     uint64          `json:"requestCounter"`
	InBytes            uint64          `json:"inBytes"`
	OutBytes           uint64          `json:"outBytes"`
	RequestMsec        uint64          `json:"requestMsec"`
	RequestMsecCounter uint64          `json:"requestMsecCounter"`
}

type statsResponses struct {
	Status1xx   uint64 `json:"1xx"`
	Status2xx   uint64 `json:"2xx"`
	Status3xx   uint64 `json:"3xx"`
	Status4xx   uint64 `json:"4xx"`
	Status5xx   uint64 `json:"5xx"`
	Miss        uint64 `json:"miss"`
	Bypass      uint64 `json:"bypass"`
	Expired     uint64 `json:"expired"`
	Stale       uint64 `json:"stale"`
	Updating    uint64 `json:"updating"`
	Revalidated uint64 `json:"revalidated"`
	Hit         uint64 `json:"hit"`
	Scarce      uint64 `json:"scarce"`
}

type statsTimeSeries struct {
	Times []int64 `json:"times"`
	Msecs []int64 `json:"msecs"`
}

type statsBuckets struct {
	Msecs    []int64 `json:"msecs"`
	Counters []int64 `json:"counters"`
}

type statsOverCounts struct {
	RequestCounter      uint64 `json:"requestCounter"`
	InBytes             uint64 `json:"inBytes"`
	OutBytes            uint64 `json:"outBytes"`
	Status1xx           uint64 `json:"1xx"`
	Status2xx           uint64 `json:"2xx"`
	Status3xx           uint64 `json:"3xx"`
	Status4xx           uint64 `json:"4xx"`
	Status5xx           uint64 `json:"5xx"`
	Miss                uint64 `json:"miss"`
	Bypass              uint64 `json:"bypass"`
	Expired             uint64 `json:"expired"`
	Stale               uint64 `json:"stale"`
	Updating            uint64 `json:"updating"`
	Revalidated         uint64 `json:"revalidated"`
	Hit                 uint64 `json:"hit"`
	Scarce              uint64 `json:"scarce"`
	RequestMsecCounter  uint64 `json:"requestMsecCounter"`
	ResponseMsecCounter uint64 `json:"responseMsecCounter"`
}

type statsUpstreamZoneData struct {
	Server              string                 `json:"server"`
	ResponseBuckets     statsBuckets           `json:"responseBuckets"`
	RequestMsecs        statsTimeSeries        `json:"requestMsecs"`
	ResponseMsecs       statsTimeSeries        `json:"responseMsecs"`
	RequestBuckets      statsBuckets           `json:"requestBuckets"`
	OverCounts          statsOverCounts        `json:"overCounts"`
	Responses           statsUpstreamResponses `json:"responses"`
	RequestMsecCounter  uint64                 `json:"requestMsecCounter"`
	ResponseMsec        uint64                 `json:"responseMsec"`
	InBytes             uint64                 `json:"inBytes"`
	RequestMsec         uint64                 `json:"requestMsec"`
	ResponseMsecCounter uint64                 `json:"responseMsecCounter"`
	RequestCounter      uint64                 `json:"requestCounter"`
	Weight              int                    `json:"weight"`
	MaxFails            int                    `json:"maxFails"`
	FailTimeout         int                    `json:"failTimeout"`
	OutBytes            uint64                 `json:"outBytes"`
	Backup              bool                   `json:"backup"`
	Down                bool                   `json:"down"`
}

type statsUpstreamResponses struct {
	Status1xx uint64 `json:"1xx"`
	Status2xx uint64 `json:"2xx"`
	Status3xx uint64 `json:"3xx"`
	Status4xx uint64 `json:"4xx"`
	Status5xx uint64 `json:"5xx"`
}

func (s *service) GetTrafficStats(ctx context.Context) (*Stats, error) {
	cfg, err := s.settingsCommands.Get(ctx)
	if err != nil {
		return nil, err
	}

	if cfg.Nginx == nil || cfg.Nginx.Stats == nil || !cfg.Nginx.Stats.Enabled {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CoreNginxStatsNotEnabled), false)
	}

	if s.semaphore.currentState() != runningState {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CoreNginxNotRunning), false)
	}

	response, err := s.fetchStatsFromSocket(ctx)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CoreNginxStatsFetchFailed), false)
	}

	return convertToStats(response), nil
}

func (s *service) fetchStatsFromSocket(ctx context.Context) (*statsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/", nil)
	if err != nil {
		return nil, err
	}

	resp, err := s.statsClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var stats statsResponse
	if err := json.Unmarshal(body, &stats); err != nil {
		return nil, err
	}

	return &stats, nil
}

func buildStatsClient(configPath string) *http.Client {
	socketPath := filepath.Join(configPath, "traffic-stats.socket")
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	return &http.Client{Transport: transport}
}

func convertToStats(src *statsResponse) *Stats {
	return &Stats{
		HostName:      src.HostName,
		Connections:   StatsConnections(src.Connections),
		ServerZones:   convertServerZones(src.ServerZones),
		FilterZones:   convertFilterZones(src.FilterZones),
		UpstreamZones: convertUpstreamZones(src.UpstreamZones),
	}
}

func convertServerZones(src map[string]statsZoneData) map[string]StatsZoneData {
	if src == nil {
		return nil
	}

	result := make(map[string]StatsZoneData, len(src))
	for k, v := range src {
		result[k] = convertZoneData(v)
	}
	return result
}

func convertFilterZones(
	src map[string]map[string]statsZoneData,
) map[string]map[string]StatsZoneData {
	if src == nil {
		return nil
	}

	result := make(map[string]map[string]StatsZoneData, len(src))
	for k, v := range src {
		if v == nil {
			result[k] = nil
			continue
		}
		inner := make(map[string]StatsZoneData, len(v))
		for ik, iv := range v {
			inner[ik] = convertZoneData(iv)
		}
		result[k] = inner
	}
	return result
}

func convertUpstreamZones(
	src map[string][]statsUpstreamZoneData,
) map[string][]StatsUpstreamZoneData {
	if src == nil {
		return nil
	}

	result := make(map[string][]StatsUpstreamZoneData, len(src))
	for k, v := range src {
		if v == nil {
			result[k] = nil
			continue
		}
		arr := make([]StatsUpstreamZoneData, len(v))
		for i, item := range v {
			arr[i] = convertUpstreamZoneData(item)
		}
		result[k] = arr
	}
	return result
}

func convertZoneData(src statsZoneData) StatsZoneData {
	return StatsZoneData{
		RequestCounter:     src.RequestCounter,
		InBytes:            src.InBytes,
		OutBytes:           src.OutBytes,
		Responses:          StatsResponses(src.Responses),
		RequestMsec:        src.RequestMsec,
		RequestMsecCounter: src.RequestMsecCounter,
		RequestMsecs:       StatsTimeSeries(src.RequestMsecs),
		RequestBuckets:     StatsBuckets(src.RequestBuckets),
		OverCounts:         StatsOverCounts(src.OverCounts),
	}
}

func convertUpstreamZoneData(src statsUpstreamZoneData) StatsUpstreamZoneData {
	return StatsUpstreamZoneData{
		Server:              src.Server,
		RequestCounter:      src.RequestCounter,
		InBytes:             src.InBytes,
		OutBytes:            src.OutBytes,
		Responses:           StatsUpstreamResponses(src.Responses),
		RequestMsec:         src.RequestMsec,
		RequestMsecCounter:  src.RequestMsecCounter,
		RequestMsecs:        StatsTimeSeries(src.RequestMsecs),
		RequestBuckets:      StatsBuckets(src.RequestBuckets),
		ResponseMsec:        src.ResponseMsec,
		ResponseMsecCounter: src.ResponseMsecCounter,
		ResponseMsecs:       StatsTimeSeries(src.ResponseMsecs),
		ResponseBuckets:     StatsBuckets(src.ResponseBuckets),
		Weight:              src.Weight,
		MaxFails:            src.MaxFails,
		FailTimeout:         src.FailTimeout,
		Backup:              src.Backup,
		Down:                src.Down,
		OverCounts:          StatsOverCounts(src.OverCounts),
	}
}
