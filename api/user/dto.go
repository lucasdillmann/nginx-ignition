package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/google/uuid"
)

type userLoginRequestDto struct {
	Username *string `json:"username" validate:"required"`
	Password *string `json:"password" validate:"required"`
}

type userLoginResponseDto struct {
	Token string `json:"token" validate:"required"`
}

type userOnboardingStatusResponseDto struct {
	Finished bool `json:"finished" validate:"required"`
}

type userPasswordUpdateRequestDto struct {
	CurrentPassword *string `json:"currentPassword" validate:"required"`
	NewPassword     *string `json:"newPassword" validate:"required"`
}

type userRequestDto struct {
	Enabled  *bool      `json:"enabled" validate:"required"`
	Name     *string    `json:"name" validate:"required"`
	Username *string    `json:"username" validate:"required"`
	Password *string    `json:"password,omitempty"`
	Role     *user.Role `json:"role" validate:"required"`
}

type userResponseDto struct {
	ID       uuid.UUID `json:"id" validate:"required"`
	Enabled  bool      `json:"enabled" validate:"required"`
	Name     string    `json:"name" validate:"required"`
	Username string    `json:"username" validate:"required"`
	Role     user.Role `json:"role" validate:"required"`
}
