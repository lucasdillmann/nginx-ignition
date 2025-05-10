package user

import (
	"context"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"github.com/google/uuid"
	"strconv"
)

const minimumUsernameLength = 3
const minimumNameLength = 3
const minimumPasswordLength = 8

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

	switch request.Role {
	case RegularRole, AdminRole:
		break
	default:
		v.delegate.Add("role", "Invalid role")
	}

	return v.delegate.Result()
}

func minimumLengthMessage(length int) string {
	return "Should have at least " + strconv.Itoa(length) + " characters"
}

func newValidator(repository Repository) *validator {
	return &validator{
		delegate:   validation.NewValidator(),
		repository: repository,
	}
}
