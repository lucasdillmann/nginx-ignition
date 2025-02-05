package cfgfiles

import (
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/host"
	"dillmann.com.br/nginx-ignition/core/settings"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/google/uuid"
	"strings"
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

func (p *hostCertificateFileProvider) provide(_ string, hosts []*host.Host) ([]output, error) {
	var bindings []*host.Binding
	for _, h := range hosts {
		bindings = append(bindings, h.Bindings...)
	}

	cgf, err := p.settingsRepository.Get()
	if err != nil {
		return nil, err
	}

	bindings = append(bindings, cgf.GlobalBindings...)

	var outputs []output
	uniqueCertIds := map[string]bool{}

	for _, binding := range bindings {
		if binding.Type == host.HttpsBindingType && binding.CertificateID != nil {
			certId := binding.CertificateID.String()
			if !uniqueCertIds[certId] {
				uniqueCertIds[certId] = true

				output, err := p.buildCertificateFile(*binding.CertificateID)
				if err != nil {
					return nil, err
				}

				outputs = append(outputs, *output)
			}
		}
	}

	return outputs, nil
}

func (p *hostCertificateFileProvider) buildCertificateFile(certificateId uuid.UUID) (*output, error) {
	cert, err := p.certificateRepository.FindByID(certificateId)
	if err != nil {
		return nil, err
	}

	certificateChain := strings.Join(cert.CertificationChain, "\n")
	mainContents := p.convertToPemEncodedString(cert.PublicKey, &cert.PrivateKey)

	certificateChain = removeEmptyLines(certificateChain)
	mainContents = removeEmptyLines(mainContents)

	contents := fmt.Sprintf("%s\n%s", certificateChain, mainContents)

	return &output{
		name:     fmt.Sprintf("certificate-%s.pem", certificateId),
		contents: contents,
	}, nil
}

func (p *hostCertificateFileProvider) convertToPemEncodedString(publicKey string, privateKey *string) string {
	publicKeyBytes, _ := base64.StdEncoding.DecodeString(publicKey)
	var buffer strings.Builder

	if strings.Contains(string(publicKeyBytes), "BEGIN CERTIFICATE") {
		buffer.WriteString(string(publicKeyBytes))
	} else {
		certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: publicKeyBytes})
		buffer.WriteString(string(certPEM))
	}

	if privateKey != nil {
		privateKeyBytes, _ := base64.StdEncoding.DecodeString(*privateKey)

		if strings.Contains(string(privateKeyBytes), "BEGIN RSA") {
			buffer.WriteString(string(privateKeyBytes))
		} else {
			keyPEM := pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: privateKeyBytes})
			buffer.WriteString(string(keyPEM))
		}
	}

	return buffer.String()
}

func removeEmptyLines(s string) string {
	s = strings.TrimPrefix(s, "\n")
	return strings.TrimSuffix(s, "\n")
}
