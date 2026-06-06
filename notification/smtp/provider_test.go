package smtp

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/notification"
)

func Test_provider(t *testing.T) {
	t.Run("ID", func(t *testing.T) {
		t.Run("returns SMTP", func(t *testing.T) {
			provider := newProvider()

			assert.Equal(t, "SMTP", provider.ID())
		})
	})

	t.Run("Name", func(t *testing.T) {
		t.Run("returns localized name", func(t *testing.T) {
			provider := newProvider()

			require.NotNil(t, provider.Name(t.Context()))
		})
	})

	t.Run("ConfigurationFields", func(t *testing.T) {
		t.Run("returns all configuration fields", func(t *testing.T) {
			provider := newProvider()

			require.Len(t, provider.ConfigurationFields(t.Context()), 8)
		})
	})

	t.Run("ImportantInstructions", func(t *testing.T) {
		t.Run("returns setup instructions", func(t *testing.T) {
			provider := newProvider()

			require.Len(t, provider.ImportantInstructions(t.Context()), 2)
		})
	})

	t.Run("Send", func(t *testing.T) {
		t.Run("delegates to mail sender", func(t *testing.T) {
			stub := &stubMailSender{}
			provider := newProvider()
			provider.sender = stub

			err := provider.Send(t.Context(), map[string]any{
				hostFieldID:   "127.0.0.1",
				portFieldID:   "1025",
				fromFieldID:   "sender@example.com",
				toFieldID:     "recipient@example.com",
				useTLSFieldID: false,
			}, notification.Deliverable{
				Title:      "Test notification",
				Summary:    "Summary text",
				OccurredAt: time.Now(),
				Category:   notification.CategoryCertificateExpiring,
			})

			require.NoError(t, err)
			assert.Equal(t, "127.0.0.1", stub.lastSettings.host)
			assert.Equal(t, 1025, stub.lastSettings.port)
			assert.Equal(t, []string{"recipient@example.com"}, stub.lastRecipients)
			assert.True(t, strings.Contains(string(stub.lastMessage), "Test notification"))
		})
	})
}
