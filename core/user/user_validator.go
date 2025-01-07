package user

import (
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"github.com/google/uuid"
	"strconv"
)

const minimumUsernameLength = 3
const minimumNameLength = 3
const minimumPasswordLength = 8

type validator struct {
	delegate *validation.ConsistencyValidator
}

func (v *validator) validate(
	repository *Repository,
	updatedState *User,
	currentState *User,
	request *SaveRequest,
	currentUserId uuid.UUID,
) error {
	if !updatedState.Enabled && currentState != nil && currentState.Id == currentUserId {
		v.delegate.Add("enabled", "You cannot disable your own user")
	}

	if request.Password == "" && currentState == nil {
		v.delegate.Add("password", validation.ValueMissingMessage)
	}

	databaseUser, _ := (*repository).FindByUsername(updatedState.Username)
	if databaseUser != nil && databaseUser.Id != updatedState.Id {
		v.delegate.Add("username", "There's already a user with the same username")
	}

	if len(updatedState.Username) < minimumUsernameLength {
		v.delegate.Add("username", minimumLengthMessage(minimumUsernameLength))
	}

	if len(updatedState.Name) < minimumNameLength {
		v.delegate.Add("name", minimumLengthMessage(minimumNameLength))
	}

	if request.Password != "" && len(request.Password) < minimumPasswordLength {
		v.delegate.Add("password", minimumLengthMessage(minimumPasswordLength))
	}

	return v.delegate.Result()
}

func minimumLengthMessage(length int) string {
	return "Should have at least " + strconv.Itoa(length) + " characters"
}

func newValidator() *validator {
	return &validator{
		delegate: &validation.ConsistencyValidator{},
	}
}
