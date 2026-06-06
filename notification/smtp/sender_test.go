package smtp

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/notification"
)

func Test_sender(t *testing.T) {
	t.Run("parseRecipients", func(t *testing.T) {
		t.Run("parses comma-separated addresses", func(t *testing.T) {
			recipients := parseRecipients("one@example.com, two@example.com ")
			assert.Equal(t, []string{"one@example.com", "two@example.com"}, recipients)
		})

		t.Run("skips empty parts", func(t *testing.T) {
			recipients := parseRecipients("one@example.com, , two@example.com")
			assert.Equal(t, []string{"one@example.com", "two@example.com"}, recipients)
		})
	})

	t.Run("parseMailSettings", func(t *testing.T) {
		t.Run("parses saved parameters", func(t *testing.T) {
			settings, recipients := parseMailSettings(map[string]any{
				hostFieldID:        "smtp.example.com",
				portFieldID:        "587",
				fromFieldID:        "sender@example.com",
				toFieldID:          "recipient@example.com",
				useStartTLSFieldID: true,
				usernameFieldID:    "user",
				passwordFieldID:    "secret",
			})
			assert.Equal(t, "smtp.example.com", settings.host)
			assert.Equal(t, 587, settings.port)
			assert.Equal(t, "sender@example.com", settings.fromAddress)
			assert.True(t, settings.useStartTLS)
			assert.False(t, settings.useTLS)
			assert.Equal(t, "user", settings.username)
			assert.Equal(t, "secret", settings.password)
			assert.Equal(t, []string{"recipient@example.com"}, recipients)
		})
	})

	t.Run("buildMessage", func(t *testing.T) {
		t.Run("builds message with all deliverable sections", func(t *testing.T) {
			message := buildMessage(notification.Deliverable{
				Title:   "Certificate expiring",
				Summary: "The certificate will expire soon.",
				Sections: []notification.DeliverableContentSection{
					{Title: new("Details"), Body: "Renew before the deadline."},
				},
				Actions: []notification.DeliverableAction{
					{Label: "Open certificate", URL: "https://example.com/certificates/1"},
				},
			}, "alerts@example.com")

			body := string(message)
			assert.Contains(t, body, "From: alerts@example.com")
			assert.Contains(t, body, "Certificate expiring")
			assert.Contains(t, body, "The certificate will expire soon.")
			assert.Contains(t, body, "Renew before the deadline.")
			assert.Contains(t, body, "Open certificate")
			assert.Contains(t, body, "https://example.com/certificates/1")
		})
	})

	t.Run("formatHTMLBody", func(t *testing.T) {
		t.Run("escapes special characters", func(t *testing.T) {
			body := formatHTMLBody(notification.Deliverable{
				Summary: "Value <script> & \"quotes\"",
			})
			assert.Contains(t, body, "&lt;script&gt;")
			assert.Contains(t, body, "&amp;")
			assert.Contains(t, body, "&quot;quotes&quot;")
		})
	})
}
