package certificate

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/pagination"
	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func newCertificate() *certificate.Certificate {
	return &certificate.Certificate{
		ID:          uuid.New(),
		ProviderID:  "custom",
		DomainNames: []string{"example.com"},
		IssuedAt:    time.Now(),
		ValidFrom:   time.Now(),
		ValidUntil:  time.Now().AddDate(0, 3, 0),
		PrivateKey:  "private-key",
		PublicKey:   "public-key",
	}
}

func newCertificatePage() *pagination.Page[certificate.Certificate] {
	return pagination.Of([]certificate.Certificate{
		*newCertificate(),
	})
}

func newIssueCertificateRequest() *issueCertificateRequest {
	return &issueCertificateRequest{
		ProviderID:  "custom",
		DomainNames: []string{"example.com"},
		Parameters: map[string]any{
			"param1": "value1",
		},
	}
}

type sseFrame struct {
	Data  map[string]any
	Event string
}

func parseSSE(t *testing.T, body []byte) []sseFrame {
	t.Helper()

	frames := make([]sseFrame, 0)
	normalized := strings.ReplaceAll(string(body), "\r\n", "\n")

	for _, block := range strings.Split(normalized, "\n\n") {
		block = strings.TrimSpace(block)
		if block == "" || strings.HasPrefix(block, ":") {
			continue
		}

		frame := sseFrame{Data: map[string]any{}}
		dataLines := make([]string, 0)

		for _, line := range strings.Split(block, "\n") {
			if name, found := strings.CutPrefix(line, "event:"); found {
				frame.Event = strings.TrimSpace(name)
			}

			if data, found := strings.CutPrefix(line, "data:"); found {
				dataLines = append(dataLines, strings.TrimSpace(data))
			}
		}

		require.NotEmpty(t, dataLines)
		require.NoError(t, json.Unmarshal([]byte(strings.Join(dataLines, "\n")), &frame.Data))
		frames = append(frames, frame)
	}

	return frames
}

func lastFrame(t *testing.T, frames []sseFrame) sseFrame {
	t.Helper()
	require.NotEmpty(t, frames)

	return frames[len(frames)-1]
}
