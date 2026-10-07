package user

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/lucasdillmann/nginx-ignition/internal/business/core/i18n"
	"github.com/lucasdillmann/nginx-ignition/internal/business/core/validation"
)

const (
	minimumUsernameLength = 3
	minimumNameLength     = 3
	minimumPasswordLength = 8
	maximumTokenNameChars = 256
)

type validator struct {
	delegate   *validation.ConsistencyValidator
	repository Repository
}

func (v *validator) validate(
	ctx context.Context,
	updatedState *User,
	currentState *User,
	request *SaveRequest,
	currentUserID *uuid.UUID,
) error {
	if !updatedState.Enabled && currentState != nil && currentUserID != nil &&
		currentState.ID == *currentUserID {
		v.delegate.Add("enabled", i18n.M(ctx, i18n.K.CoreUserCannotDisableSelf))
	}

	if request.Password == nil && currentState == nil {
		v.delegate.Add("password", i18n.M(ctx, i18n.K.CommonValueMissing))
	}

	databaseUser, _ := v.repository.FindByUsername(ctx, updatedState.Username)
	if databaseUser != nil && databaseUser.ID != updatedState.ID {
		v.delegate.Add("username", i18n.M(ctx, i18n.K.CoreUserDuplicatedUsername))
	}

	if len(updatedState.Username) < minimumUsernameLength {
		v.delegate.Add(
			"username",
			i18n.M(ctx, i18n.K.CoreUserTooShort).V("min", minimumUsernameLength),
		)
	}

	if len(updatedState.Name) < minimumNameLength {
		v.delegate.Add(
			"name",
			i18n.M(ctx, i18n.K.CoreUserTooShort).V("min", minimumNameLength),
		)
	}

	if request.Password != nil && len(*request.Password) < minimumPasswordLength {
		v.delegate.Add(
			"password",
			i18n.M(ctx, i18n.K.CoreUserTooShort).V("min", minimumPasswordLength),
		)
	}

	v.validatePermissions(ctx, request.Permissions)

	return v.delegate.Result()
}

func (v *validator) validatePermissions(ctx context.Context, permissions Permissions) {
	v.validatePermission(ctx, "hosts", permissions.Hosts)
	v.validatePermission(ctx, "streams", permissions.Streams)
	v.validatePermission(ctx, "certificates", permissions.Certificates)
	v.validatePermission(ctx, "logs", permissions.Logs)
	v.validatePermission(ctx, "integrations", permissions.Integrations)
	v.validatePermission(ctx, "accessLists", permissions.AccessLists)
	v.validatePermission(ctx, "settings", permissions.Settings)
	v.validatePermission(ctx, "users", permissions.Users)
	v.validatePermission(ctx, "nginxServer", permissions.NginxServer)
	v.validatePermission(ctx, "exportData", permissions.ExportData)
	v.validatePermission(ctx, "vpns", permissions.VPNs)
	v.validatePermission(ctx, "caches", permissions.Caches)
	v.validatePermission(ctx, "trafficStats", permissions.TrafficStats)

	if permissions.NginxServer == NoAccessAccessLevel {
		v.delegate.Add("permissions.nginxServer", i18n.M(ctx, i18n.K.CoreUserAtLeastReadOnly))
	}

	if permissions.Logs == ReadWriteAccessLevel {
		v.delegate.Add("permissions.logs", i18n.M(ctx, i18n.K.CoreUserCannotHaveWriteAccess))
	}

	if permissions.ExportData == ReadWriteAccessLevel {
		v.delegate.Add(
			"permissions.exportData",
			i18n.M(ctx, i18n.K.CoreUserCannotHaveWriteAccess),
		)
	}

	if permissions.TrafficStats == ReadWriteAccessLevel {
		v.delegate.Add(
			"permissions.trafficStats",
			i18n.M(ctx, i18n.K.CoreUserCannotHaveWriteAccess),
		)
	}
}

func (v *validator) validatePermission(ctx context.Context, key string, value AccessLevel) {
	switch value {
	case NoAccessAccessLevel, ReadOnlyAccessLevel, ReadWriteAccessLevel:
	default:
		v.delegate.Add(
			fmt.Sprintf("permissions.%s", key),
			i18n.M(ctx, i18n.K.CoreUserInvalidAccessLevel),
		)
	}
}

func newValidator(repository Repository) *validator {
	return &validator{
		delegate:   validation.NewValidator(),
		repository: repository,
	}
}

func (v *validator) validateAPIToken(
	ctx context.Context,
	userID uuid.UUID,
	request *NewAPITokenRequest,
) error {
	name := strings.TrimSpace(request.Name)

	if name == "" {
		v.delegate.Add("name", i18n.M(ctx, i18n.K.CommonCannotBeEmpty))
	} else {
		exists, err := v.repository.ExistsTokenByName(ctx, userID, name)
		if err != nil {
			return err
		}

		if exists {
			v.delegate.Add("name", i18n.M(ctx, i18n.K.CoreUserDuplicatedTokenName))
		}
	}

	if len(name) > maximumTokenNameChars {
		v.delegate.Add(
			"name",
			i18n.M(ctx, i18n.K.CommonValueTooLong).V("max", maximumTokenNameChars),
		)
	}

	if request.Expiration != nil && !request.Expiration.After(time.Now()) {
		v.delegate.Add("expiration", i18n.M(ctx, i18n.K.CoreUserTokenExpiredDate))
	}

	return v.delegate.Result()
}
