package nginx

func newMetadata() *Metadata {
	return &Metadata{
		Version:       "1.25.3",
		BuildDetails:  "built by gcc",
		Modules:       []string{},
		tlsSniEnabled: true,
	}
}
