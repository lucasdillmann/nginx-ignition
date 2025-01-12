package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
)

func toDomain(model *userModel) *user.User {
	return &user.User{
		ID:           model.ID,
		Enabled:      model.Enabled,
		Name:         model.Name,
		Username:     model.Username,
		PasswordHash: model.PasswordHash,
		PasswordSalt: model.PasswordSalt,
		Role:         user.Role(model.Role),
	}
}

func toModel(domain *user.User) *userModel {
	return &userModel{
		ID:           domain.ID,
		Enabled:      domain.Enabled,
		Name:         domain.Name,
		Username:     domain.Username,
		PasswordHash: domain.PasswordHash,
		PasswordSalt: domain.PasswordSalt,
		Role:         string(domain.Role),
	}
}
