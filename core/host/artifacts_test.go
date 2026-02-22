package host

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/certificate"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

type validatorMocks struct {
	repository  *MockedRepository
	integration *integration.MockedCommands
	vpn         *vpn.MockedCommands
	accessList  *accesslist.MockedCommands
	cache       *cache.MockedCommands
	binding     *binding.MockedCommands
	certificate *certificate.MockedCommands
}

func (m *validatorMocks) newValidator() *validator {
	return newValidator(
		m.repository,
		m.integration,
		m.vpn,
		m.accessList,
		m.cache,
		m.binding,
		m.certificate,
	)
}

func newHost() *Host {
	return &Host{
		ID:          uuid.New(),
		Enabled:     true,
		DomainNames: []string{"example.com"},
		Bindings: []binding.Binding{
			{
				Type: binding.HTTPBindingType,
				IP:   "0.0.0.0",
				Port: 80,
			},
		},
		Routes: []Route{
			{
				ID:         uuid.New(),
				Enabled:    true,
				Priority:   0,
				SourcePath: "/",
				Type:       StaticResponseRouteType,
				Response: &RouteStaticResponse{
					StatusCode: 200,
					Payload:    new("OK"),
				},
			},
		},
		FeatureSet: FeatureSet{
			StatsEnabled: true,
		},
	}
}

func setupValidator(t *testing.T) (*validator, *validatorMocks) {
	ctrl := gomock.NewController(t)
	repo := NewMockedRepository(ctrl)
	integrationCmds := integration.NewMockedCommands(ctrl)
	vpnCmds := vpn.NewMockedCommands(ctrl)
	aclCmds := accesslist.NewMockedCommands(ctrl)
	cacheCmds := cache.NewMockedCommands(ctrl)
	bindingCmds := binding.NewMockedCommands(ctrl)
	certCmds := certificate.NewMockedCommands(ctrl)

	mocks := &validatorMocks{
		repository:  repo,
		integration: integrationCmds,
		vpn:         vpnCmds,
		accessList:  aclCmds,
		cache:       cacheCmds,
		binding:     bindingCmds,
		certificate: certCmds,
	}

	return mocks.newValidator(), mocks
}

func assertViolations(t *testing.T, err error, msgs ...string) {
	t.Helper()

	if assert.Error(t, err) {
		var consistencyErr *validation.ConsistencyError
		if assert.ErrorAs(t, err, &consistencyErr) {
			for _, msg := range msgs {
				found := false
				for _, v := range consistencyErr.Violations {
					if strings.Contains(v.Message.Key, msg) {
						found = true
						break
					}
				}
				if !found {
					allMsgs := make([]string, 0, len(consistencyErr.Violations))
					for _, v := range consistencyErr.Violations {
						allMsgs = append(allMsgs, fmt.Sprintf("'%s'", v.Message.String()))
					}
					assert.Failf(
						t,
						"Violation not found",
						"Expected violation containing '%s', got: [%s]",
						msg,
						strings.Join(allMsgs, ", "),
					)
				}
			}
		} else {
			assert.Failf(
				t,
				"Unexpected error type",
				"Expected ConsistencyError, got %T: %v",
				err,
				err,
			)
		}
	}
}
