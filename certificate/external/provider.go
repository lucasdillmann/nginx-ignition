package external

import (
	"context"
	"encoding/base64"
	"time"

	"github.com/google/uuid"

	"dillmann.com.br/nginx-ignition/certificate/commons"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/dynamicfields"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
)

type Provider struct{}

func newProvider() *Provider {
	return &Provider{}
}

func (p *Provider) ID() string {
	return providerID
}

func (p *Provider) Name(ctx context.Context) *i18n.Message {
	return i18n.M(ctx, i18n.K.CertificateExternalName)
}

func (p *Provider) DynamicFields(ctx context.Context) []dynamicfields.DynamicField {
	return dynamicFields(ctx)
}

func (p *Provider) Priority() int {
	return 3
}

func (p *Provider) Issue(
	ctx context.Context,
	request *certificate.IssueRequest,
) (*certificate.Certificate, error) {
	fields := p.DynamicFields(ctx)
	if err := commons.Validate(ctx, request, &validationRules{fields}); err != nil {
		return nil, err
	}

	return loadCertificateFromParameters(
		ctx,
		request.Parameters,
		request.DomainNames,
		uuid.New(),
	)
}

func (p *Provider) Renew(
	ctx context.Context,
	existing *certificate.Certificate,
) (*certificate.Certificate, error) {
	fields := p.DynamicFields(ctx)
	renewRequest := &certificate.IssueRequest{
		Parameters:  existing.Parameters,
		ProviderID:  p.ID(),
		DomainNames: existing.DomainNames,
	}

	if err := commons.Validate(ctx, renewRequest, &validationRules{fields}); err != nil {
		return nil, err
	}

	return loadCertificateFromParameters(
		ctx,
		existing.Parameters,
		existing.DomainNames,
		existing.ID,
	)
}

func (p *Provider) IsDueToRenew(
	ctx context.Context,
	existing *certificate.Certificate,
) (bool, error) {
	if existing.RenewAfter != nil {
		return !time.Now().Before(*existing.RenewAfter), nil
	}

	publicPath, found := existing.Parameters[publicKeyPathFieldID].(string)
	if !found || publicPath == "" {
		return false, nil
	}

	publicCert, err := readCertificate(ctx, publicPath)
	if err != nil {
		return false, err
	}

	onDisk := base64.StdEncoding.EncodeToString(publicCert.Raw)
	if existing.PublicKey != onDisk {
		return true, nil
	}

	return false, nil
}
