package nginx

import (
	"net/http"
)

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func newMetadata() *Metadata {
	return &Metadata{
		Version:       "1.25.3",
		BuildDetails:  "built by gcc",
		Modules:       []string{},
		tlsSniEnabled: true,
	}
}
