package netbird

type noOpLogger struct{}

func (n *noOpLogger) Write(_ []byte) (int, error) {
	// NO-OP
	return 0, nil
}
