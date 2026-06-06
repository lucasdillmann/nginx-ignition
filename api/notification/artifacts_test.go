package notification

import (
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/notification"
)

func sampleNotification(userID uuid.UUID) *notification.Notification {
	notificationID := uuid.New()

	return &notification.Notification{
		ID:        notificationID,
		UserID:    userID,
		Title:     "Certificate renewed",
		Summary:   "example.com was renewed",
		Category:  notification.CategoryCertificateRenewed,
		ReadAt:    new(time.Now()),
		CreatedAt: time.Now(),
		Payload: notification.Payload{
			OccurredAt: time.Now(),
			Tags:       map[string]string{"domain": "example.com"},
			Sections: []notification.DeliverableContentSection{
				{Body: "Renewal completed"},
			},
			Actions: []notification.DeliverableAction{
				{Label: "View", URL: "/certificates"},
			},
		},
		RelatedEntities: []notification.RelatedEntity{
			{Type: "certificate", ID: uuid.New(), Name: "example.com"},
		},
		Submissions: []notification.ProviderSubmission{
			{Provider: "SMTP", Status: notification.ProviderSubmissionStatusSuccess},
		},
		DeliveryCompleted: true,
	}
}
