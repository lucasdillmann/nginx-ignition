package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/common/validation"
)

const (
	minimumUsernameLength = 3
	minimumNameLength     = 3
	minimumPasswordLength = 8
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
	currentUserId *uuid.UUID,
) error {
	if !updatedState.Enabled && currentState != nil && currentUserId != nil && currentState.ID == *currentUserId {
		v.delegate.Add("enabled", "You cannot disable your own user")
	}

	if request.Password == nil && currentState == nil {
		v.delegate.Add("password", validation.ValueMissingMessage)
	}

	databaseUser, _ := v.repository.FindByUsername(ctx, updatedState.Username)
	if databaseUser != nil && databaseUser.ID != updatedState.ID {
		v.delegate.Add("username", "There's already a user with the same username")
	}

	if len(updatedState.Username) < minimumUsernameLength {
		v.delegate.Add("username", minimumLengthMessage(minimumUsernameLength))
	}

	if len(updatedState.Name) < minimumNameLength {
		v.delegate.Add("name", minimumLengthMessage(minimumNameLength))
	}

	if request.Password != nil && len(*request.Password) < minimumPasswordLength {
		v.delegate.Add("password", minimumLengthMessage(minimumPasswordLength))
	}

	v.validatePermissions(request.Permissions)

	return v.delegate.Result()
}

func (v *validator) validatePermissions(permissions Permissions) {
	v.validatePermission("hosts", permissions.Hosts)
	v.validatePermission("streams", permissions.Streams)
	v.validatePermission("certificates", permissions.Certificates)
	v.validatePermission("logs", permissions.Logs)
	v.validatePermission("integrations", permissions.Integrations)
	v.validatePermission("accessLists", permissions.AccessLists)
	v.validatePermission("settings", permissions.Settings)
	v.validatePermission("users", permissions.Users)
	v.validatePermission("nginxServer", permissions.NginxServer)
	v.validatePermission("exportData", permissions.ExportData)
	v.validatePermission("vpns", permissions.VPNs)
	v.validatePermission("caches", permissions.Caches)

	if permissions.NginxServer == NoAccessAccessLevel {
		v.delegate.Add("permissions.nginxServer", "At least read-only access is required")
	}

	if permissions.Logs == ReadWriteAccessLevel {
		v.delegate.Add("permissions.logs", "Cannot have read-write access to logs")
	}

	if permissions.ExportData == ReadWriteAccessLevel {
		v.delegate.Add("permissions.exportData", "Cannot have read-write access to data export")
	}
}

func (v *validator) validatePermission(key string, value AccessLevel) {
	switch value {
	case NoAccessAccessLevel, ReadOnlyAccessLevel, ReadWriteAccessLevel:
	default:
		v.delegate.Add(fmt.Sprintf("permissions.%s", key), "Invalid access level")
	}
}

func minimumLengthMessage(length int) string {
	return fmt.Sprintf("Should have at least %d characters", length)
}

func newValidator(repository Repository) *validator {
	return &validator{
		delegate:   validation.NewValidator(),
		repository: repository,
	}
}
