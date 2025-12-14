package client

import (
	"dillmann.com.br/nginx-ignition/core/certificate/client"
)

func toDomain(model *certificateModel, clients []*clientModel) (*client.Certificate, error) {
	var ca *client.CA
	if model.CAPublicKey != nil || model.CAPrivateKey != nil || model.CASendToClients != nil {
		ca = &client.CA{
			PublicKey:     model.CAPublicKey,
			PrivateKey:    model.CAPrivateKey,
			SendToClients: model.CASendToClients,
		}
	}

	stapling := client.Stapling{
		Enabled:           model.StaplingEnabled,
		Verify:            model.StaplingVerify,
		ResponderURL:      model.StaplingResponderURL,
		ResponderFilePath: model.StaplingResponderFilePath,
	}

	domainClients := make([]*client.Client, 0, len(clients))
	for _, clientModel := range clients {
		domainClients = append(domainClients, &client.Client{
			ID:         clientModel.ID,
			DN:         clientModel.DN,
			PublicKey:  clientModel.PublicKey,
			PrivateKey: clientModel.PrivateKey,
			IssuedAt:   clientModel.IssuedAt,
			ExpiresAt:  clientModel.ExpiresAt,
			Revoked:    clientModel.Revoked,
		})
	}

	return &client.Certificate{
		ID:             model.ID,
		Name:           model.Name,
		Type:           client.Type(model.Type),
		CA:             ca,
		Clients:        domainClients,
		ValidationMode: client.ValidationMode(model.ValidationMode),
		Stapling:       stapling,
	}, nil
}

func toModel(certificate *client.Certificate) (*certificateModel, []*clientModel, error) {
	certModel := &certificateModel{
		ID:                        certificate.ID,
		Name:                      certificate.Name,
		Type:                      string(certificate.Type),
		ValidationMode:            string(certificate.ValidationMode),
		StaplingEnabled:           certificate.Stapling.Enabled,
		StaplingVerify:            certificate.Stapling.Verify,
		StaplingResponderURL:      certificate.Stapling.ResponderURL,
		StaplingResponderFilePath: certificate.Stapling.ResponderFilePath,
	}

	if certificate.CA != nil {
		certModel.CAPublicKey = certificate.CA.PublicKey
		certModel.CAPrivateKey = certificate.CA.PrivateKey
		certModel.CASendToClients = certificate.CA.SendToClients
	}

	var clientModels []*clientModel
	for _, clientCert := range certificate.Clients {
		clientModels = append(clientModels, &clientModel{
			ID:                  clientCert.ID,
			ClientCertificateID: certificate.ID,
			DN:                  clientCert.DN,
			PublicKey:           clientCert.PublicKey,
			PrivateKey:          clientCert.PrivateKey,
			IssuedAt:            clientCert.IssuedAt,
			ExpiresAt:           clientCert.ExpiresAt,
			Revoked:             clientCert.Revoked,
		})
	}

	return certModel, clientModels, nil
}
