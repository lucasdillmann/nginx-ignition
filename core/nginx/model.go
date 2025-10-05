package nginx

type Metadata struct {
	Version       string
	BuildDetails  string
	TLSSNIEnabled bool
	Modules       []string
}

func (m *Metadata) StreamSupportType() string {
	if m.hasModule("stream") {
		return "static"
	} else if m.hasModule("ngx_stream_module") {
		return "dynamic"
	} else {
		return "none"
	}
}

func (m *Metadata) RunCodeSupportAvailable() bool {
	jsModule := m.hasModule("ngx_http_js_module")
	luaModule := m.hasModule("ngx_http_lua_module")
	ndkModule := m.hasModule("ndk_http_module")

	return jsModule && luaModule && ndkModule
}

func (m *Metadata) hasModule(name string) bool {
	for _, module := range m.Modules {
		if module == name {
			return true
		}
	}

	return false
}
