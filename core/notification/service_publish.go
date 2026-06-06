package notification

import (
	"context"
	"time"

	"github.com/google/uuid"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/coreerror"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/log"
)

func (s *service) Publish(
	ctx context.Context,
	userID uuid.UUID,
	request SendRequest,
) (*Notification, error) {
	return s.publishForUser(ctx, userID, request)
}

func (s *service) Broadcast(ctx context.Context, request SendRequest) error {
	userIDs, err := s.userCommands.ListEnabledIDs(ctx)
	if err != nil {
		return err
	}

	for _, userID := range userIDs {
		if _, err := s.publishForUser(ctx, userID, request); err != nil {
			log.Errorf("notification broadcast failed for user %s: %s", userID, err)
		}
	}

	return nil
}

func (s *service) publishForUser(
	ctx context.Context,
	userID uuid.UUID,
	request SendRequest,
) (*Notification, error) {
	if !isValidCategory(request.Category) {
		return nil, coreerror.New(
			i18n.M(ctx, i18n.K.CoreNotificationInvalidCategory),
			true,
		)
	}

	usr, err := s.userCommands.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if usr == nil {
		return nil, coreerror.New(i18n.M(ctx, i18n.K.CoreUserNotFound), false)
	}

	notificationLanguage := usr.NotificationLanguage
	lang := resolveLanguage(s.i18nCommands, notificationLanguage)
	deliverable := resolveSendRequest(s.i18nCommands, lang, request)
	notification := notificationFromDeliverable(userID, deliverable)

	relatedEntities := make([]StoredRelatedEntity, len(request.RelatedEntities))
	for index, entity := range request.RelatedEntities {
		relatedEntities[index] = StoredRelatedEntity{
			NotificationID: notification.ID,
			Type:           entity.Type,
			ID:             entity.ID,
			Name:           entity.Name,
		}
	}

	if saveErr := s.repository.SaveNotification(
		ctx,
		notification,
		relatedEntities,
	); saveErr != nil {
		return nil, saveErr
	}

	configurations, err := s.repository.FindEnabledConfigurationsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	submissions := make([]ProviderSubmission, 0)
	for _, configuration := range configurations {
		if !configurationAcceptsCategory(configuration.Categories, request.Category) {
			continue
		}

		submissions = append(submissions, ProviderSubmission{
			ID:              uuid.New(),
			NotificationID:  notification.ID,
			ConfigurationID: configuration.ID,
			Provider:        configuration.Provider,
			Status:          SubmissionStatusPending,
			AttemptCount:    0,
		})
	}

	if len(submissions) > 0 {
		if err := s.repository.SaveProviderSubmissions(ctx, submissions); err != nil {
			return nil, err
		}
	} else {
		if err := s.repository.SetDeliveryCompleted(ctx, notification.ID, true); err != nil {
			return nil, err
		}
		notification.DeliveryCompleted = true
	}

	notification.RelatedEntities = toRelatedEntities(relatedEntities)
	return notification, nil
}

func resolveSendRequest(
	commands i18n.Commands,
	lang language.Tag,
	request SendRequest,
) Deliverable {
	sections := make([]DeliverableContentSection, len(request.Sections))
	for index, section := range request.Sections {
		resolved := DeliverableContentSection{
			Body: translateDetached(commands, lang, section.Body),
		}
		if section.Title != nil {
			resolved.Title = new(translateDetached(commands, lang, *section.Title))
		}
		sections[index] = resolved
	}

	actions := make([]DeliverableAction, len(request.Actions))
	for index, action := range request.Actions {
		actions[index] = DeliverableAction{
			Label: translateDetached(commands, lang, action.Label),
			URL:   action.URL,
		}
	}

	tags := request.Tags
	if tags == nil {
		tags = map[string]string{}
	}

	return Deliverable{
		Title:      translateDetached(commands, lang, request.Title),
		Summary:    translateDetached(commands, lang, request.Summary),
		Sections:   sections,
		Actions:    actions,
		OccurredAt: request.OccurredAt,
		Tags:       tags,
		Category:   request.Category,
	}
}

func translateDetached(
	commands i18n.Commands,
	lang language.Tag,
	detached i18n.DetachedMessage,
) string {
	variables := detached.Variables
	if variables == nil {
		variables = map[string]any{}
	}

	return commands.Translate(lang, detached.Key, variables)
}

func resolveLanguage(commands i18n.Commands, languageValue string) language.Tag {
	if languageValue == "" {
		return commands.DefaultLanguage()
	}

	tag, err := language.Parse(languageValue)
	if err != nil || !commands.Supports(tag) {
		return commands.DefaultLanguage()
	}

	return tag
}

func notificationFromDeliverable(userID uuid.UUID, deliverable Deliverable) *Notification {
	now := time.Now()
	occurredAt := deliverable.OccurredAt
	if occurredAt.IsZero() {
		occurredAt = now
	}

	tags := deliverable.Tags
	if tags == nil {
		tags = map[string]string{}
	}

	return &Notification{
		ID:        uuid.New(),
		UserID:    userID,
		Title:     deliverable.Title,
		Summary:   deliverable.Summary,
		Category:  deliverable.Category,
		CreatedAt: now,
		Payload: Payload{
			Sections:   deliverable.Sections,
			Actions:    deliverable.Actions,
			OccurredAt: occurredAt,
			Tags:       tags,
		},
	}
}
