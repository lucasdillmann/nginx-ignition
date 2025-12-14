package client

import (
	"dillmann.com.br/nginx-ignition/core/certificate/client"
)

func toResponseDto(input *client.Certificate) *certificateResponseDto {
	if input == nil {
		return nil
	}

	response := &certificateResponseDto{
		ID:             input.ID,
		Name:           input.Name,
		Type:           string(input.Type),
		ValidationMode: string(input.ValidationMode),
		Stapling: &staplingResponseDto{
			Enabled:           input.Stapling.Enabled,
			Verify:            input.Stapling.Verify,
			ResponderURL:      input.Stapling.ResponderURL,
			ResponderFilePath: input.Stapling.ResponderFilePath,
		},
	}

	if input.CA != nil && input.CA.SendToClients != nil {
		response.SendCAToClients = input.CA.SendToClients
	}

	if len(input.Clients) > 0 {
		clients := make([]clientResponseDto, len(input.Clients))
		for index, item := range input.Clients {
			clients[index] = clientResponseDto{
				ID:        item.ID,
				DN:        item.DN,
				IssuedAt:  item.IssuedAt,
				ExpiresAt: item.ExpiresAt,
				Revoked:   item.Revoked,
			}
		}

		response.Clients = &clients
	}

	return response
}
