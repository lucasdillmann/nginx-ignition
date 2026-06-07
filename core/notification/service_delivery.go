package notification

import (
	"context"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (s *service) ProcessPendingDeliveries(ctx context.Context) error {
	notifications, err := s.repository.FindNotificationsWithIncompleteDelivery(ctx)
	if err != nil {
		return err
	}

	for _, notification := range notifications {
		if err := s.processNotificationDeliveries(ctx, &notification); err != nil {
			log.Errorf(
				"notification delivery failed for notification %s: %s",
				notification.ID,
				err,
			)
		}
	}

	return nil
}

func (s *service) processNotificationDeliveries(
	ctx context.Context,
	notification *Notification,
) error {
	submissions, err := s.repository.FindPendingSubmissionsByNotificationID(
		ctx,
		notification.ID,
	)
	if err != nil {
		return err
	}

	deliverable := Deliverable{
		Title:      notification.Title,
		Summary:    notification.Summary,
		Sections:   notification.Payload.Sections,
		Actions:    notification.Payload.Actions,
		OccurredAt: notification.Payload.OccurredAt,
		Tags:       notification.Payload.Tags,
		Category:   notification.Category,
	}

	for index := range submissions {
		if err := s.processPendingSubmission(
			ctx,
			notification,
			&submissions[index],
			deliverable,
		); err != nil {
			return err
		}
	}

	return s.updateDeliveryCompletedIfNeeded(ctx, notification.ID)
}

func (s *service) processPendingSubmission(
	ctx context.Context,
	notification *Notification,
	submission *ProviderSubmission,
	deliverable Deliverable,
) error {
	provider := s.findProvider(submission.Provider)
	if provider == nil {
		submission.Status = SubmissionStatusFailed
		submission.LastAttemptAt = new(time.Now())
		submission.LastError = new("provider not found")
		submission.AttemptCount++

		log.Errorf(submissionDeliveryFailureLogFormat, submission.ID, "provider not found")
		if err := s.repository.UpdateProviderSubmission(ctx, submission); err != nil {
			log.Errorf(submissionDeliveryFailureLogFormat, submission.ID, err)
		}

		return nil
	}

	configuration, err := s.repository.FindConfigurationByIDAndUserID(
		ctx,
		submission.ConfigurationID,
		notification.UserID,
	)
	if err != nil {
		return err
	}

	if configuration == nil {
		submission.Status = SubmissionStatusFailed
		submission.LastAttemptAt = new(time.Now())
		submission.LastError = new("configuration not found")
		submission.AttemptCount++

		log.Errorf(submissionDeliveryFailureLogFormat, submission.ID, "configuration not found")
		if err := s.repository.UpdateProviderSubmission(ctx, submission); err != nil {
			log.Errorf(submissionDeliveryFailureLogFormat, submission.ID, err)
		}

		return nil
	}

	if !configuration.Enabled {
		submission.Status = SubmissionStatusSkipped
		submission.SucceededAt = nil
		submission.LastError = nil

		if err := s.repository.UpdateProviderSubmission(ctx, submission); err != nil {
			log.Warnf(
				"notification delivery skipped for submission %s: %s",
				submission.ID,
				err,
			)
		}

		return nil
	}

	s.attemptProviderDelivery(ctx, provider, configuration, submission, deliverable)
	return nil
}

func (s *service) attemptProviderDelivery(
	ctx context.Context,
	provider Provider,
	configuration *Configuration,
	submission *ProviderSubmission,
	deliverable Deliverable,
) {
	sendError := provider.Send(ctx, configuration.Parameters, deliverable)
	submission.LastAttemptAt = new(time.Now())
	submission.AttemptCount++

	if sendError == nil {
		submission.Status = SubmissionStatusSuccess
		submission.SucceededAt = submission.LastAttemptAt
		submission.LastError = nil
	} else {
		submission.LastError = new(sendError.Error())

		log.Errorf(submissionDeliveryFailureLogFormat, submission.ID, sendError)
		if submission.AttemptCount >= maxDeliveryAttempts {
			submission.Status = SubmissionStatusFailed
		}
	}

	if err := s.repository.UpdateProviderSubmission(ctx, submission); err != nil {
		log.Errorf(submissionDeliveryFailureLogFormat, submission.ID, err)
	}
}

func (s *service) updateDeliveryCompletedIfNeeded(
	ctx context.Context,
	notificationID uuid.UUID,
) error {
	submissions, err := s.repository.FindSubmissionsByNotificationID(ctx, notificationID)
	if err != nil {
		return err
	}

	if len(submissions) == 0 {
		return s.repository.SetDeliveryCompleted(ctx, notificationID, true)
	}

	for _, submission := range submissions {
		if !submissionStatusIsTerminal(submission.Status) {
			return nil
		}
	}

	return s.repository.SetDeliveryCompleted(ctx, notificationID, true)
}
