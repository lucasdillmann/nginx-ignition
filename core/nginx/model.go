package nginx

type SupportType string

const (
	StaticSupportType  SupportType = "STATIC"
	DynamicSupportType SupportType = "DYNAMIC"
	NoneSupportType    SupportType = "NONE"
)

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
