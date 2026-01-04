package host

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/accesslist"
	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/cache"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/common/validation"
	"dillmann.com.br/nginx-ignition/core/integration"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_validator(t *testing.T) {
	ctx := context.Background()

	t.Run("validate", func(t *testing.T) {
		t.Run("valid host passes", func(t *testing.T) {
			hostValidator, _, _, _, _, _, bindingCmds := setupValidator(t)
			h := newHost()

			bindingCmds.EXPECT().
				Validate(ctx, "bindings", 0, &h.Bindings[0], gomock.Any()).
				Return(nil)

			assert.NoError(t, hostValidator.validate(ctx, h))
		})

		t.Run("validates simple host fields", func(t *testing.T) {
			t.Run("invalid domain names", func(t *testing.T) {
				hostValidator, _, _, _, _, _, bindingCmds := setupValidator(t)
				h := newHost()
				h.DomainNames = []string{"invalid_domain"}

				bindingCmds.EXPECT().
					Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(ctx, h)
				assertViolations(t, err, "Value is not a valid domain name")
			})

			t.Run("default server logic", func(t *testing.T) {
				t.Run("error if another default exists", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, _ := setupValidator(t)
					h := newHost()
					h.DefaultServer = true
					h.DomainNames = nil
					h.Bindings = nil

					otherDefault := newHost()
					otherDefault.ID = uuid.New()
					repo.EXPECT().FindDefault(ctx).Return(otherDefault, nil)

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "already another host marked as the default")
				})

				t.Run("error if domains provided for default", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, _ := setupValidator(t)
					h := newHost()
					h.DefaultServer = true
					h.DomainNames = []string{"example.com"}
					h.Bindings = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil)

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "Must be empty when the host is the default one")
				})
			})
		})

		t.Run("validates bindings", func(t *testing.T) {
			t.Run("global bindings logic", func(t *testing.T) {
				hostValidator, _, _, _, _, _, _ := setupValidator(t)
				h := newHost()
				h.UseGlobalBindings = true
				h.Bindings = []binding.Binding{{}}

				err := hostValidator.validate(ctx, h)
				assertViolations(t, err, "Must be empty when using global bindings")
			})

			t.Run("custom bindings required", func(t *testing.T) {
				hostValidator, _, _, _, _, _, _ := setupValidator(t)
				h := newHost()
				h.UseGlobalBindings = false
				h.Bindings = nil

				err := hostValidator.validate(ctx, h)
				assertViolations(t, err, "At least one binding must be informed")
			})
		})

		t.Run("validates routes", func(t *testing.T) {
			t.Run("requires at least one route", func(t *testing.T) {
				hostValidator, _, _, _, _, _, bindingCmds := setupValidator(t)
				h := newHost()
				h.Routes = nil

				bindingCmds.EXPECT().
					Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(ctx, h)
				assertViolations(t, err, "At least one route must be informed")
			})

			t.Run("duplicates priority", func(t *testing.T) {
				hostValidator, _, _, _, _, _, bindingCmds := setupValidator(t)
				h := newHost()
				h.Routes = []Route{
					{
						Priority:   10,
						SourcePath: "/a",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: ptr.Of("ok")},
					},
					{
						Priority:   10,
						SourcePath: "/b",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: ptr.Of("ok")},
					},
				}

				bindingCmds.EXPECT().
					Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(ctx, h)
				assertViolations(t, err, "Priority 10 is duplicated")
			})

			t.Run("duplicates source path", func(t *testing.T) {
				hostValidator, _, _, _, _, _, bindingCmds := setupValidator(t)
				h := newHost()
				h.Routes = []Route{
					{
						Priority:   10,
						SourcePath: "/a",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: ptr.Of("ok")},
					},
					{
						Priority:   20,
						SourcePath: "/a",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: ptr.Of("ok")},
					},
				}

				bindingCmds.EXPECT().
					Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(ctx, h)
				assertViolations(t, err, "Source path was already used")
			})

			t.Run("validates route types", func(t *testing.T) {
				t.Run("Proxy", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, bindingCmds := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = ProxyRouteType
					h.Routes[0].TargetURI = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
					bindingCmds.EXPECT().
						Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is required")

					h.Routes[0].TargetURI = ptr.Of("http://invalid\nurl")
					hostValidator = newValidator(repo, nil, nil, nil, nil, bindingCmds)
					err = hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is not a valid URL")
				})

				t.Run("Redirect", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, bindingCmds := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = RedirectRouteType
					h.Routes[0].TargetURI = nil
					h.Routes[0].RedirectCode = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
					bindingCmds.EXPECT().
						Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is required")

					h.Routes[0].TargetURI = ptr.Of("http://example.com")
					h.Routes[0].RedirectCode = ptr.Of(200)
					hostValidator = newValidator(repo, nil, nil, nil, nil, bindingCmds)
					err = hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value must be between 300 and 399")
				})

				t.Run("StaticResponse", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, bindingCmds := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = StaticResponseRouteType
					h.Routes[0].Response = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
					bindingCmds.EXPECT().
						Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "A value is required")

					h.Routes[0].Response = &RouteStaticResponse{StatusCode: 999}
					hostValidator = newValidator(repo, nil, nil, nil, nil, bindingCmds)
					err = hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value must be between 100 and 599")
				})

				t.Run("Integration", func(t *testing.T) {
					hostValidator, repo, integrationCmds, _, _, _, bindingCmds := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = IntegrationRouteType
					h.Routes[0].Integration = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
					bindingCmds.EXPECT().
						Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is required")

					integrationID := uuid.New()
					h.Routes[0].Integration = &RouteIntegrationConfig{
						IntegrationID: integrationID,
						OptionID:      "",
					}
					integrationCmds.EXPECT().Exists(ctx, integrationID).Return(ptr.Of(false), nil)

					hostValidator = newValidator(repo, integrationCmds, nil, nil, nil, bindingCmds)
					err = hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is required")
				})

				t.Run("ExecuteCode", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, bindingCmds := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = ExecuteCodeRouteType
					h.Routes[0].SourceCode = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
					bindingCmds.EXPECT().
						Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is required")

					h.Routes[0].SourceCode = &RouteSourceCode{Language: "INVALID", Contents: ""}
					hostValidator = newValidator(repo, nil, nil, nil, nil, bindingCmds)
					err = hostValidator.validate(ctx, h)
					assertViolations(t, err, "Invalid value", "Value is required")
				})

				t.Run("StaticFiles", func(t *testing.T) {
					hostValidator, repo, _, _, _, _, bindingCmds := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = StaticFilesRouteType
					h.Routes[0].TargetURI = nil

					repo.EXPECT().FindDefault(ctx).Return(nil, nil).AnyTimes()
					bindingCmds.EXPECT().
						Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value is required")

					h.Routes[0].TargetURI = ptr.Of("invalid/path")
					hostValidator = newValidator(repo, nil, nil, nil, nil, bindingCmds)
					err = hostValidator.validate(ctx, h)
					assertViolations(t, err, "Value must start with a /")
				})
			})
		})

		t.Run("validates VPNs", func(t *testing.T) {
			hostValidator, _, _, vpnCmds, _, _, bindingCmds := setupValidator(t)
			h := newHost()
			vpnID := uuid.New()
			h.VPNs = []VPN{{VPNID: vpnID, Name: "vpn1"}}

			bindingCmds.EXPECT().
				Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			vpnCmds.EXPECT().Get(ctx, vpnID).Return(nil, nil)

			err := hostValidator.validate(ctx, h)
			assertViolations(t, err, "No VPN connection was found")
		})

		t.Run("validates ACLs", func(t *testing.T) {
			hostValidator, _, _, _, aclCmds, _, bindingCmds := setupValidator(t)
			h := newHost()
			aclID := uuid.New()
			h.AccessListID = &aclID

			bindingCmds.EXPECT().
				Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			aclCmds.EXPECT().Exists(ctx, aclID).Return(false, nil)

			err := hostValidator.validate(ctx, h)
			assertViolations(t, err, "No access list found")
		})

		t.Run("validates Cache", func(t *testing.T) {
			hostValidator, _, _, _, _, cacheCmds, bindingCmds := setupValidator(t)
			h := newHost()
			cacheID := uuid.New()
			h.CacheID = &cacheID

			bindingCmds.EXPECT().
				Validate(ctx, gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			cacheCmds.EXPECT().Exists(ctx, cacheID).Return(false, nil)

			err := hostValidator.validate(ctx, h)
			assertViolations(t, err, "No cache configuration found")
		})
	})
}

func setupValidator(t *testing.T) (
	*validator,
	*MockedRepository,
	*integration.MockedCommands,
	*vpn.MockedCommands,
	*accesslist.MockedCommands,
	*cache.MockedCommands,
	*binding.MockedCommands,
) {
	ctrl := gomock.NewController(t)
	repo := NewMockedRepository(ctrl)
	integrationCmds := integration.NewMockedCommands(ctrl)
	vpnCmds := vpn.NewMockedCommands(ctrl)
	aclCmds := accesslist.NewMockedCommands(ctrl)
	cacheCmds := cache.NewMockedCommands(ctrl)
	bindingCmds := binding.NewMockedCommands(ctrl)
	hostValidator := newValidator(repo, integrationCmds, vpnCmds, aclCmds, cacheCmds, bindingCmds)

	return hostValidator, repo, integrationCmds, vpnCmds, aclCmds, cacheCmds, bindingCmds
}

func assertViolations(t *testing.T, err error, msgs ...string) {
	t.Helper()

	if assert.Error(t, err) {
		var consistencyErr *validation.ConsistencyError
		if assert.ErrorAs(t, err, &consistencyErr) {
			for _, msg := range msgs {
				found := false
				for _, v := range consistencyErr.Violations {
					if strings.Contains(v.Message, msg) {
						found = true
						break
					}
				}
				if !found {
					var allMsgs []string
					for _, v := range consistencyErr.Violations {
						allMsgs = append(allMsgs, fmt.Sprintf("'%s'", v.Message))
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
			assert.Failf(t, "Unexpected error type", "Expected ConsistencyError, got %T: %v", err, err)
		}
	}
}
