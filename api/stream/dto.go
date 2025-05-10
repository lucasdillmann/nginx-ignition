package stream

import (
	"dillmann.com.br/nginx-ignition/core/stream"
	"github.com/google/uuid"
)

type streamRequestDto struct {
	Enabled     *bool          `json:"enabled" validate:"required"`
	Description *string        `json:"description" validate:"required"`
	FeatureSet  *featureSetDto `json:"featureSet" validate:"required"`
	Backend     *addressDto    `json:"backend" validate:"required"`
	Binding     *addressDto    `json:"binding" validate:"required"`
}

type featureSetDto struct {
	UseProxyProtocol *bool `json:"useProxyProtocol"`
	SSL              *bool `json:"ssl"`
	TCPKeepAlive     *bool `json:"tcpKeepAlive"`
	TCPNoDelay       *bool `json:"tcpNoDelay"`
	TCPDeferred      *bool `json:"tcpDeferred"`
}

type addressDto struct {
	Protocol stream.Protocol `json:"protocol"`
	Address  *string         `json:"address"`
	Port     *int            `json:"port"`
}

type streamResponseDto struct {
	ID          *uuid.UUID     `json:"id" validate:"required"`
	Enabled     *bool          `json:"enabled" validate:"required"`
	Description *string        `json:"description" validate:"required"`
	FeatureSet  *featureSetDto `json:"featureSet" validate:"required"`
	Backend     *addressDto    `json:"backend" validate:"required"`
	Binding     *addressDto    `json:"binding" validate:"required"`
}
