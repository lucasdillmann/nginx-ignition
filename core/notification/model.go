package notification

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Category string

const (
	CategoryCertificateRenewed     Category = "CERTIFICATE_RENEWED"
	CategoryCertificateRenewFailed Category = "CERTIFICATE_RENEW_FAILED"
	CategoryCertificateExpiring    Category = "CERTIFICATE_EXPIRING"
	CategoryNginxReloadFailed      Category = "NGINX_RELOAD_FAILED"
	CategoryNginxReloadSucceeded   Category = "NGINX_RELOAD_SUCCEEDED"
)

func AllCategories() []Category {
	return []Category{
		CategoryCertificateRenewed,
		CategoryCertificateRenewFailed,
		CategoryCertificateExpiring,
		CategoryNginxReloadFailed,
		CategoryNginxReloadSucceeded,
	}
}

func CategoryName(ctx context.Context, category Category) *i18n.Message {
	switch category {
	case CategoryCertificateRenewed:
		return i18n.M(ctx, i18n.K.CoreNotificationCategoryCertificateRenewed)
	case CategoryCertificateRenewFailed:
		return i18n.M(ctx, i18n.K.CoreNotificationCategoryCertificateRenewFailed)
	case CategoryCertificateExpiring:
		return i18n.M(ctx, i18n.K.CoreNotificationCategoryCertificateExpiring)
	case CategoryNginxReloadFailed:
		return i18n.M(ctx, i18n.K.CoreNotificationCategoryNginxReloadFailed)
	case CategoryNginxReloadSucceeded:
		return i18n.M(ctx, i18n.K.CoreNotificationCategoryNginxReloadSucceeded)
	default:
		return i18n.M(ctx, i18n.K.CoreNotificationCategoryUnknown)
	}
}

func isValidCategory(category Category) bool {
	switch category {
	case CategoryCertificateRenewed,
		CategoryCertificateRenewFailed,
		CategoryCertificateExpiring,
		CategoryNginxReloadFailed,
		CategoryNginxReloadSucceeded:
		return true
	default:
		return false
	}
}

func configurationAcceptsCategory(allowed *[]Category, category Category) bool {
	if allowed == nil {
		return true
	}

	if len(*allowed) == 0 {
		return false
	}

	for _, item := range *allowed {
		if item == category {
			return true
		}
	}

	return false
}

type CategoryInfo struct {
	Name *i18n.Message
	ID   Category
}

type SubmissionStatus string

const (
	SubmissionStatusPending SubmissionStatus = "PENDING"
	SubmissionStatusSuccess SubmissionStatus = "SUCCESS"
	SubmissionStatusFailed  SubmissionStatus = "FAILED"
	SubmissionStatusSkipped SubmissionStatus = "SKIPPED"
)

type ProviderSubmissionStatus = SubmissionStatus

const (
	ProviderSubmissionStatusPending = SubmissionStatusPending
	ProviderSubmissionStatusSuccess = SubmissionStatusSuccess
	ProviderSubmissionStatusFailed  = SubmissionStatusFailed
	ProviderSubmissionStatusSkipped = SubmissionStatusSkipped
)

func submissionStatusIsTerminal(status SubmissionStatus) bool {
	switch status {
	case SubmissionStatusSuccess, SubmissionStatusFailed, SubmissionStatusSkipped:
		return true
	default:
		return false
	}
}

type Payload struct {
	OccurredAt time.Time
	Tags       map[string]string
	Sections   []DeliverableContentSection
	Actions    []DeliverableAction
}

type DeliverableContentSection struct {
	Title *string
	Body  string
}

type DeliverableAction struct {
	Label string
	URL   string
}

type Notification struct {
	Payload           Payload
	CreatedAt         time.Time
	ReadAt            *time.Time
	Title             string
	Summary           string
	Category          Category
	RelatedEntities   []RelatedEntity
	Submissions       []ProviderSubmission
	UserID            uuid.UUID
	ID                uuid.UUID
	DeliveryCompleted bool
}

type StoredRelatedEntity struct {
	Name           string
	Type           string
	NotificationID uuid.UUID
	ID             uuid.UUID
}

type ProviderSubmission struct {
	LastError       *string
	LastAttemptAt   *time.Time
	SucceededAt     *time.Time
	Provider        string
	Status          SubmissionStatus
	NotificationID  uuid.UUID
	ConfigurationID uuid.UUID
	ID              uuid.UUID
	AttemptCount    int
}

type Configuration struct {
	Categories *[]Category
	Parameters map[string]any
	Name       string
	Provider   string
	UserID     uuid.UUID
	ID         uuid.UUID
	Enabled    bool
}
