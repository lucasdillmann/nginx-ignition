package cfgfiles

import (
	"context"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
)

type hostCertificateFileProvider struct {
	certificateRepository certificate.Repository
	settingsRepository    settings.Repository
}

func newHostCertificateFileProvider(certificateRepository certificate.Repository, settingsRepository settings.Repository) *hostCertificateFileProvider {
	return &hostCertificateFileProvider{
		certificateRepository: certificateRepository,
		settingsRepository:    settingsRepository,
	}
}

func (p *hostCertificateFileProvider) provide(ctx *providerContext) ([]File, error) {
	var bindings []host.Binding
	for _, h := range ctx.hosts {
		bindings = append(bindings, h.Bindings...)
	}

	cgf, err := p.settingsRepository.Get(ctx.context)
	if err != nil {
		return nil, err
	}

	bindings = append(bindings, cgf.GlobalBindings...)

	var outputs []File
	uniqueCertIds := map[string]bool{}

	for _, binding := range bindings {
		if binding.Type == host.HttpsBindingType && binding.CertificateID != nil {
			certId := binding.CertificateID.String()
			if !uniqueCertIds[certId] {
				uniqueCertIds[certId] = true

				output, err := p.buildCertificateFile(ctx.context, *binding.CertificateID)
				if err != nil {
					return nil, err
				}

				outputs = append(outputs, *output)
			}
		}
	}

	return outputs, nil
}

func (p *hostCertificateFileProvider) buildCertificateFile(
	ctx context.Context,
	certificateId uuid.UUID,
) (*File, error) {
	cert, err := p.certificateRepository.FindByID(ctx, certificateId)
	if err != nil {
		return nil, err
	}

	certBytes, _ := base64.StdEncoding.DecodeString(cert.PublicKey)
	encodedCert := convertToPemEncodedCertificateString(certBytes)

	privateKeyBytes, _ := base64.StdEncoding.DecodeString(cert.PrivateKey)
	encodedPrivateKey := convertToPemEncodedPrivateKeyString(privateKeyBytes)

	var certificateChain string
	for _, chainElement := range cert.CertificationChain {
		decodedBytes, _ := base64.StdEncoding.DecodeString(chainElement)
		certificateChain += convertToPemEncodedCertificateString(decodedBytes) + "\n"
	}

	var contents string
	if len(certificateChain) > 0 {
		contents = fmt.Sprintf("%s\n%s\n%s", encodedCert, certificateChain, encodedPrivateKey)
	} else {
		contents = fmt.Sprintf("%s\n%s", encodedCert, encodedPrivateKey)
	}

	return &File{
		Name:     fmt.Sprintf("certificate-%s.pem", certificateId),
		Contents: contents,
	}, nil
}

func convertToPemEncodedCertificateString(bytes []byte) string {
	if strings.Contains(string(bytes), "CERTIFICATE") {
		return string(bytes)
	} else {
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: bytes})
		return string(certPEM)
	}
}

func convertToPemEncodedPrivateKeyString(bytes []byte) string {
	if strings.Contains(string(bytes), "PRIVATE KEY") {
		return string(bytes)
	} else {
		keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: bytes})
		return string(keyPEM)
	}
}
