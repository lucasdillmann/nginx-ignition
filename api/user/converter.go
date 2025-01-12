package user

import (
	"dillmann.com.br/nginx-ignition/core/user"
	"github.com/google/uuid"
)

func toDomain(dto *userRequestDto) *user.SaveRequest {
	return &user.SaveRequest{
		ID:       uuid.New(),
		Enabled:  *dto.Enabled,
		Name:     *dto.Name,
		Username: *dto.Username,
		Password: dto.Password,
		Role:     *dto.Role,
	}
}

func toDto(domain *user.User) *userResponseDto {
	return &userResponseDto{
		ID:       domain.ID,
		Enabled:  domain.Enabled,
		Name:     domain.Name,
		Username: domain.Username,
		Role:     domain.Role,
	}
}
