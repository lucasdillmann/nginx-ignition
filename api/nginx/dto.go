package nginx

type trafficStatsResponseDTO struct {
	ServerZones   map[string]trafficStatsZoneDataDTO            `json:"serverZones"`
	FilterZones   map[string]map[string]trafficStatsZoneDataDTO `json:"filterZones"`
	UpstreamZones map[string][]trafficStatsUpstreamZoneDataDTO  `json:"upstreamZones"`
	HostName      string                                        `json:"hostName"`
	Connections   trafficStatsConnectionsDTO                    `json:"connections"`
}

type trafficStatsConnectionsDTO struct {
	Active   uint64 `json:"active"`
	Reading  uint64 `json:"reading"`
	Writing  uint64 `json:"writing"`
	Waiting  uint64 `json:"waiting"`
	Accepted uint64 `json:"accepted"`
	Handled  uint64 `json:"handled"`
	Requests uint64 `json:"requests"`
}

type trafficStatsZoneDataDTO struct {
	RequestMsecs       trafficStatsTimeSeriesDTO `json:"requestMsecs"`
	Responses          trafficStatsResponsesDTO  `json:"responses"`
	RequestCounter     uint64                    `json:"requestCounter"`
	InBytes            uint64                    `json:"inBytes"`
	OutBytes           uint64                    `json:"outBytes"`
	RequestMsec        uint64                    `json:"requestMsec"`
	RequestMsecCounter uint64                    `json:"requestMsecCounter"`
}

type trafficStatsResponsesDTO struct {
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

type trafficStatsTimeSeriesDTO struct {
	Times []int64 `json:"times"`
	Msecs []int64 `json:"msecs"`
}

type trafficStatsUpstreamZoneDataDTO struct {
	Server              string                           `json:"server"`
	RequestMsecs        trafficStatsTimeSeriesDTO        `json:"requestMsecs"`
	ResponseMsecs       trafficStatsTimeSeriesDTO        `json:"responseMsecs"`
	Responses           trafficStatsUpstreamResponsesDTO `json:"responses"`
	RequestMsecCounter  uint64                           `json:"requestMsecCounter"`
	ResponseMsec        uint64                           `json:"responseMsec"`
	InBytes             uint64                           `json:"inBytes"`
	RequestMsec         uint64                           `json:"requestMsec"`
	ResponseMsecCounter uint64                           `json:"responseMsecCounter"`
	RequestCounter      uint64                           `json:"requestCounter"`
	Weight              int                              `json:"weight"`
	MaxFails            int                              `json:"maxFails"`
	FailTimeout         int                              `json:"failTimeout"`
	OutBytes            uint64                           `json:"outBytes"`
	Backup              bool                             `json:"backup"`
	Down                bool                             `json:"down"`
}

type trafficStatsUpstreamResponsesDTO struct {
	Status1xx uint64 `json:"1xx"`
	Status2xx uint64 `json:"2xx"`
	Status3xx uint64 `json:"3xx"`
	Status4xx uint64 `json:"4xx"`
	Status5xx uint64 `json:"5xx"`
}
