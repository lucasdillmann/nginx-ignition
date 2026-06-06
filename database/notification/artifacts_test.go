package notification

import (
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

func newNotification(userID uuid.UUID) *notification.Notification {
	return &notification.Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     "Certificate renewed",
		Summary:   "Certificate example.com was renewed",
		Category:  notification.CategoryCertificateRenewed,
		CreatedAt: time.Now().UTC().Truncate(time.Second),
		Payload: notification.Payload{
			Sections: []notification.DeliverableContentSection{
				{Body: "Renewal completed successfully."},
			},
			Actions: []notification.DeliverableAction{
				{Label: "View certificate", URL: "/certificates"},
			},
			OccurredAt: time.Now().UTC().Truncate(time.Second),
			Tags:       map[string]string{"domain": "example.com"},
		},
		DeliveryCompleted: false,
	}
}

func newConfiguration(userID uuid.UUID) *notification.Configuration {
	return &notification.Configuration{
		ID:       uuid.New(),
		UserID:   userID,
		Name:     uuid.NewString(),
		Provider: "SMTP",
		Enabled:  true,
		Parameters: map[string]any{
			"host": "smtp.example.com",
		},
	}
}
