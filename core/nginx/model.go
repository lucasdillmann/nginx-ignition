package nginx

type SupportType string

const (
	StaticSupportType  SupportType = "STATIC"
	DynamicSupportType SupportType = "DYNAMIC"
	NoneSupportType    SupportType = "NONE"
)

type Stats struct {
	ServerZones   map[string]StatsZoneData
	FilterZones   map[string]map[string]StatsZoneData
	UpstreamZones map[string][]StatsUpstreamZoneData
	HostName      string
	Connections   StatsConnections
}

type StatsConnections struct {
	Active   uint64
	Reading  uint64
	Writing  uint64
	Waiting  uint64
	Accepted uint64
	Handled  uint64
	Requests uint64
}

type StatsZoneData struct {
	RequestMsecs       StatsTimeSeries
	Responses          StatsResponses
	RequestCounter     uint64
	InBytes            uint64
	OutBytes           uint64
	RequestMsec        uint64
	RequestMsecCounter uint64
}

type StatsResponses struct {
	Status1xx   uint64
	Status2xx   uint64
	Status3xx   uint64
	Status4xx   uint64
	Status5xx   uint64
	Miss        uint64
	Bypass      uint64
	Expired     uint64
	Stale       uint64
	Updating    uint64
	Revalidated uint64
	Hit         uint64
	Scarce      uint64
}

type StatsTimeSeries struct {
	Times []int64
	Msecs []int64
}

type StatsUpstreamZoneData struct {
	Server              string
	RequestMsecs        StatsTimeSeries
	ResponseMsecs       StatsTimeSeries
	Responses           StatsUpstreamResponses
	RequestMsecCounter  uint64
	ResponseMsec        uint64
	InBytes             uint64
	RequestMsec         uint64
	ResponseMsecCounter uint64
	RequestCounter      uint64
	Weight              int
	MaxFails            int
	FailTimeout         int
	OutBytes            uint64
	Backup              bool
	Down                bool
}

type StatsUpstreamResponses struct {
	Status1xx uint64
	Status2xx uint64
	Status3xx uint64
	Status4xx uint64
	Status5xx uint64
}

type Metadata struct {
	Version       string
	BuildDetails  string
	Modules       []string
	tlsSniEnabled bool
}

func (m *Metadata) SNISupportType() SupportType {
	if m.tlsSniEnabled {
		return StaticSupportType
	}

	return NoneSupportType
}

func (m *Metadata) StreamSupportType() SupportType {
	if m.hasModule("stream") {
		return StaticSupportType
	}

	if m.hasModule("ngx_stream_module") {
		return DynamicSupportType
	}

	return NoneSupportType
}

func (m *Metadata) StatsSupportType() SupportType {
	if m.hasModule("nginx-module-vts") || m.hasModule("ngx_http_vts_module") {
		return DynamicSupportType
	}

	return NoneSupportType
}

func (m *Metadata) RunCodeSupportType() SupportType {
	jsModule := m.hasModule("ngx_http_js_module")
	luaModule := m.hasModule("ngx_http_lua_module")
	ndkModule := m.hasModule("ndk_http_module")

	if jsModule && luaModule && ndkModule {
		return DynamicSupportType
	}

	return NoneSupportType
}

func (m *Metadata) hasModule(name string) bool {
	for _, module := range m.Modules {
		if module == name {
			return true
		}
	}

	return false
}
