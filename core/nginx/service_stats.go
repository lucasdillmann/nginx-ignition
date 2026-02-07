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
	HostName      string                              `json:"hostName"`
	Connections   statsConnections                    `json:"connections"`
	ServerZones   map[string]statsZoneData            `json:"serverZones"`
	FilterZones   map[string]map[string]statsZoneData `json:"filterZones"`
	UpstreamZones map[string][]statsUpstreamZoneData  `json:"upstreamZones"`
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
	RequestCounter     uint64          `json:"requestCounter"`
	InBytes            uint64          `json:"inBytes"`
	OutBytes           uint64          `json:"outBytes"`
	Responses          statsResponses  `json:"responses"`
	RequestMsec        uint64          `json:"requestMsec"`
	RequestMsecCounter uint64          `json:"requestMsecCounter"`
	RequestMsecs       statsTimeSeries `json:"requestMsecs"`
	RequestBuckets     statsBuckets    `json:"requestBuckets"`
	OverCounts         statsOverCounts `json:"overCounts"`
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
	RequestCounter      uint64                 `json:"requestCounter"`
	InBytes             uint64                 `json:"inBytes"`
	OutBytes            uint64                 `json:"outBytes"`
	Responses           statsUpstreamResponses `json:"responses"`
	RequestMsec         uint64                 `json:"requestMsec"`
	RequestMsecCounter  uint64                 `json:"requestMsecCounter"`
	RequestMsecs        statsTimeSeries        `json:"requestMsecs"`
	RequestBuckets      statsBuckets           `json:"requestBuckets"`
	ResponseMsec        uint64                 `json:"responseMsec"`
	ResponseMsecCounter uint64                 `json:"responseMsecCounter"`
	ResponseMsecs       statsTimeSeries        `json:"responseMsecs"`
	ResponseBuckets     statsBuckets           `json:"responseBuckets"`
	Weight              int                    `json:"weight"`
	MaxFails            int                    `json:"maxFails"`
	FailTimeout         int                    `json:"failTimeout"`
	Backup              bool                   `json:"backup"`
	Down                bool                   `json:"down"`
	OverCounts          statsOverCounts        `json:"overCounts"`
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

	socketPath := filepath.Join(s.processManager.configPath, "traffic-stats.socket")
	response, err := s.fetchStatsFromSocket(ctx, socketPath)
	if err != nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CoreNginxStatsFetchFailed), false)
	}

	return convertToStats(response), nil
}

func (s *service) fetchStatsFromSocket(ctx context.Context, socketPath string) (*statsResponse, error) {
	transport := &http.Transport{
		DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
			return net.Dial("unix", socketPath)
		},
	}

	client := &http.Client{Transport: transport}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost/", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
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

func convertToStats(src *statsResponse) *Stats {
	return &Stats{
		HostName:      src.HostName,
		Connections:   convertConnections(src.Connections),
		ServerZones:   convertServerZones(src.ServerZones),
		FilterZones:   convertFilterZones(src.FilterZones),
		UpstreamZones: convertUpstreamZones(src.UpstreamZones),
	}
}

func convertConnections(src statsConnections) StatsConnections {
	return StatsConnections{
		Active:   src.Active,
		Reading:  src.Reading,
		Writing:  src.Writing,
		Waiting:  src.Waiting,
		Accepted: src.Accepted,
		Handled:  src.Handled,
		Requests: src.Requests,
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

func convertFilterZones(src map[string]map[string]statsZoneData) map[string]map[string]StatsZoneData {
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

func convertUpstreamZones(src map[string][]statsUpstreamZoneData) map[string][]StatsUpstreamZoneData {
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
		Responses:          convertResponses(src.Responses),
		RequestMsec:        src.RequestMsec,
		RequestMsecCounter: src.RequestMsecCounter,
		RequestMsecs:       convertTimeSeries(src.RequestMsecs),
		RequestBuckets:     convertBuckets(src.RequestBuckets),
		OverCounts:         convertOverCounts(src.OverCounts),
	}
}

func convertResponses(src statsResponses) StatsResponses {
	return StatsResponses{
		Status1xx:   src.Status1xx,
		Status2xx:   src.Status2xx,
		Status3xx:   src.Status3xx,
		Status4xx:   src.Status4xx,
		Status5xx:   src.Status5xx,
		Miss:        src.Miss,
		Bypass:      src.Bypass,
		Expired:     src.Expired,
		Stale:       src.Stale,
		Updating:    src.Updating,
		Revalidated: src.Revalidated,
		Hit:         src.Hit,
		Scarce:      src.Scarce,
	}
}

func convertTimeSeries(src statsTimeSeries) StatsTimeSeries {
	return StatsTimeSeries{
		Times: src.Times,
		Msecs: src.Msecs,
	}
}

func convertBuckets(src statsBuckets) StatsBuckets {
	return StatsBuckets{
		Msecs:    src.Msecs,
		Counters: src.Counters,
	}
}

func convertOverCounts(src statsOverCounts) StatsOverCounts {
	return StatsOverCounts{
		RequestCounter:      src.RequestCounter,
		InBytes:             src.InBytes,
		OutBytes:            src.OutBytes,
		Status1xx:           src.Status1xx,
		Status2xx:           src.Status2xx,
		Status3xx:           src.Status3xx,
		Status4xx:           src.Status4xx,
		Status5xx:           src.Status5xx,
		Miss:                src.Miss,
		Bypass:              src.Bypass,
		Expired:             src.Expired,
		Stale:               src.Stale,
		Updating:            src.Updating,
		Revalidated:         src.Revalidated,
		Hit:                 src.Hit,
		Scarce:              src.Scarce,
		RequestMsecCounter:  src.RequestMsecCounter,
		ResponseMsecCounter: src.ResponseMsecCounter,
	}
}

func convertUpstreamZoneData(src statsUpstreamZoneData) StatsUpstreamZoneData {
	return StatsUpstreamZoneData{
		Server:              src.Server,
		RequestCounter:      src.RequestCounter,
		InBytes:             src.InBytes,
		OutBytes:            src.OutBytes,
		Responses:           convertUpstreamResponses(src.Responses),
		RequestMsec:         src.RequestMsec,
		RequestMsecCounter:  src.RequestMsecCounter,
		RequestMsecs:        convertTimeSeries(src.RequestMsecs),
		RequestBuckets:      convertBuckets(src.RequestBuckets),
		ResponseMsec:        src.ResponseMsec,
		ResponseMsecCounter: src.ResponseMsecCounter,
		ResponseMsecs:       convertTimeSeries(src.ResponseMsecs),
		ResponseBuckets:     convertBuckets(src.ResponseBuckets),
		Weight:              src.Weight,
		MaxFails:            src.MaxFails,
		FailTimeout:         src.FailTimeout,
		Backup:              src.Backup,
		Down:                src.Down,
		OverCounts:          convertOverCounts(src.OverCounts),
	}
}

func convertUpstreamResponses(src statsUpstreamResponses) StatsUpstreamResponses {
	return StatsUpstreamResponses{
		Status1xx: src.Status1xx,
		Status2xx: src.Status2xx,
		Status3xx: src.Status3xx,
		Status4xx: src.Status4xx,
		Status5xx: src.Status5xx,
	}
}
