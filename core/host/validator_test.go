package host

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"dillmann.com.br/nginx-ignition/core/binding"
	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/vpn"
)

func Test_validator(t *testing.T) {
	t.Run("validate", func(t *testing.T) {
		t.Run("valid host passes", func(t *testing.T) {
			hostValidator, mocks := setupValidator(t)
			h := newHost()

			mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
			mocks.binding.EXPECT().
				Validate(t.Context(), "bindings", 0, &h.Bindings[0], gomock.Any()).
				Return(nil)

			assert.NoError(t, hostValidator.validate(t.Context(), h))
		})

		t.Run("validates simple host fields", func(t *testing.T) {
			t.Run("invalid domain names", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				h.DomainNames = []string{"invalid_domain"}

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
				mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
				mocks.binding.EXPECT().
					Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CommonInvalidDomainName)
			})

			t.Run("default server logic", func(t *testing.T) {
				t.Run("error if another default exists", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.DefaultServer = true
					h.DomainNames = nil
					h.Bindings = nil

					otherDefault := newHost()
					otherDefault.ID = uuid.New()
					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(otherDefault, nil)

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostDefaultAlreadyExists)
				})

				t.Run("error if domains provided for default", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.DefaultServer = true
					h.DomainNames = []string{"example.com"}
					h.Bindings = nil

					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil)
					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostDomainMustBeEmptyForDefault)
				})
			})
		})

		t.Run("validates bindings", func(t *testing.T) {
			t.Run("global bindings logic", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				h.UseGlobalBindings = true
				h.Bindings = []binding.Binding{{}}

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CoreHostBindingsMustBeEmptyForGlobal)
			})

			t.Run("custom bindings required", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				h.UseGlobalBindings = false
				h.Bindings = nil

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CommonAtLeastOneRequired)
			})
		})

		t.Run("validates routes", func(t *testing.T) {
			t.Run("requires at least one route", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				h.Routes = nil

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
				mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
				mocks.binding.EXPECT().
					Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CommonAtLeastOneRequired)
			})

			t.Run("duplicates priority", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				h.Routes = []Route{
					{
						Priority:   10,
						SourcePath: "/a",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: new("ok")},
					},
					{
						Priority:   10,
						SourcePath: "/b",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: new("ok")},
					},
				}

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
				mocks.binding.EXPECT().
					Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CoreHostDuplicatedRoutePriority)
			})

			t.Run("duplicates source path", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				h.Routes = []Route{
					{
						Priority:   10,
						SourcePath: "/a",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: new("ok")},
					},
					{
						Priority:   20,
						SourcePath: "/a",
						Type:       StaticResponseRouteType,
						Response:   &RouteStaticResponse{StatusCode: 200, Payload: new("ok")},
					},
				}

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
				mocks.binding.EXPECT().
					Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CoreHostDuplicatedSourcePath)
			})

			t.Run("validates route types", func(t *testing.T) {
				t.Run("Proxy", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = ProxyRouteType
					h.Routes[0].TargetURI = nil

					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostTargetUriRequired)

					h.Routes[0].TargetURI = new("http://invalid\nurl")
					hostValidator = mocks.newValidator()
					err = hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CommonInvalidUrl)
				})

				t.Run("Redirect", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = RedirectRouteType
					h.Routes[0].TargetURI = nil
					h.Routes[0].RedirectCode = nil

					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostTargetUriRequired)

					h.Routes[0].TargetURI = new("http://example.com")
					h.Routes[0].RedirectCode = new(200)
					hostValidator = mocks.newValidator()
					err = hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CommonBetweenValues)
				})

				t.Run("StaticResponse", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = StaticResponseRouteType
					h.Routes[0].Response = nil

					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostStaticResponseRequired)

					h.Routes[0].Response = &RouteStaticResponse{StatusCode: 999}
					hostValidator = mocks.newValidator()
					err = hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CommonBetweenValues)
				})

				t.Run("Integration", func(t *testing.T) {
					t.Run("missing integration data", func(t *testing.T) {
						hostValidator, mocks := setupValidator(t)
						h := newHost()
						h.Routes[0].Type = IntegrationRouteType
						h.Routes[0].Integration = nil

						mocks.vpn.EXPECT().
							GetAvailableDrivers(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.repository.EXPECT().
							FindDefault(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.binding.EXPECT().
							Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(nil).
							AnyTimes()

						err := hostValidator.validate(t.Context(), h)
						assertViolations(t, err, i18n.K.CoreHostIntegrationRequired)
					})

					t.Run("missing integration option", func(t *testing.T) {
						hostValidator, mocks := setupValidator(t)
						h := newHost()
						h.Routes[0].Type = IntegrationRouteType
						integrationID := uuid.New()
						h.Routes[0].Integration = &RouteIntegrationConfig{
							IntegrationID: integrationID,
							OptionID:      "",
						}

						mocks.vpn.EXPECT().
							GetAvailableDrivers(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.repository.EXPECT().
							FindDefault(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.binding.EXPECT().
							Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(nil).
							AnyTimes()
						mocks.integration.EXPECT().
							Exists(t.Context(), integrationID).
							Return(new(false), nil)

						err := hostValidator.validate(t.Context(), h)
						assertViolations(t, err, i18n.K.CoreHostIntegrationRequired)
					})

					t.Run("invalid TargetURI", func(t *testing.T) {
						hostValidator, mocks := setupValidator(t)
						h := newHost()
						h.Routes[0].Type = IntegrationRouteType
						h.Routes[0].Integration = &RouteIntegrationConfig{
							IntegrationID: uuid.New(),
							OptionID:      "opt-1",
						}
						h.Routes[0].TargetURI = new("invalid uri")

						mocks.vpn.EXPECT().
							GetAvailableDrivers(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.repository.EXPECT().
							FindDefault(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.binding.EXPECT().
							Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(nil).
							AnyTimes()
						mocks.integration.EXPECT().
							Exists(t.Context(), h.Routes[0].Integration.IntegrationID).
							Return(new(true), nil)

						err := hostValidator.validate(t.Context(), h)
						assertViolations(t, err, i18n.K.CoreHostInvalidUri)
					})

					t.Run("valid TargetURI", func(t *testing.T) {
						hostValidator, mocks := setupValidator(t)
						h := newHost()
						h.Routes[0].Type = IntegrationRouteType
						h.Routes[0].Integration = &RouteIntegrationConfig{
							IntegrationID: uuid.New(),
							OptionID:      "opt-1",
						}
						h.Routes[0].TargetURI = new("/api/v1")

						mocks.vpn.EXPECT().
							GetAvailableDrivers(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.repository.EXPECT().
							FindDefault(t.Context()).
							Return(nil, nil).
							AnyTimes()
						mocks.binding.EXPECT().
							Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
							Return(nil).
							AnyTimes()
						mocks.integration.EXPECT().
							Exists(t.Context(), h.Routes[0].Integration.IntegrationID).
							Return(new(true), nil)

						err := hostValidator.validate(t.Context(), h)
						assert.NoError(t, err)
					})
				})

				t.Run("ExecuteCode", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = ExecuteCodeRouteType
					h.Routes[0].SourceCode = nil

					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostSourceCodeRequired)

					h.Routes[0].SourceCode = &RouteSourceCode{Language: "INVALID", Contents: ""}
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()
					err = hostValidator.validate(t.Context(), h)
					assertViolations(
						t,
						err,
						i18n.K.CommonInvalidValue,
						i18n.K.CoreHostSourceCodeRequired,
					)
				})

				t.Run("StaticFiles", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					h.Routes[0].Type = StaticFilesRouteType
					h.Routes[0].TargetURI = nil

					mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostTargetUriRequired)

					h.Routes[0].TargetURI = new("invalid/path")
					mocks.repository.EXPECT().FindDefault(t.Context()).Return(nil, nil).AnyTimes()
					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).
						AnyTimes()
					err = hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CommonStartsWithSlashRequired)
				})
			})
		})

		t.Run("validates VPNs", func(t *testing.T) {
			t.Run("vpn not found", func(t *testing.T) {
				hostValidator, mocks := setupValidator(t)
				h := newHost()
				vpnID := uuid.New()
				h.VPNs = []VPN{{VPNID: vpnID, Name: "vpn1"}}

				mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
				mocks.binding.EXPECT().
					Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(nil)
				mocks.vpn.EXPECT().
					GetAvailableDrivers(t.Context()).
					Return(nil, nil).
					AnyTimes().
					AnyTimes()
				mocks.vpn.EXPECT().Get(t.Context(), vpnID).Return(nil, nil)

				err := hostValidator.validate(t.Context(), h)
				assertViolations(t, err, i18n.K.CoreHostVpnNotFound)
			})

			t.Run("vpn certificate validation", func(t *testing.T) {
				t.Run("certificate informed but https disabled", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					vpnID := uuid.New()
					certID := uuid.New()
					h.VPNs = []VPN{
						{VPNID: vpnID, Name: "vpn1", EnableHTTPS: false, CertificateID: &certID},
					}

					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).AnyTimes()
					mocks.vpn.EXPECT().
						Get(t.Context(), vpnID).
						Return(&vpn.VPN{Enabled: true, Driver: "driver1"}, nil)
					mocks.vpn.EXPECT().
						GetAvailableDrivers(t.Context()).
						Return(nil, nil).AnyTimes()

					err := hostValidator.validate(t.Context(), h)
					assertViolations(
						t,
						err,
						i18n.K.CoreHostVpnCertificateCannotBeInformedIfDisabled,
					)
				})

				t.Run("driver managed - certificate required", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					vpnID := uuid.New()
					h.VPNs = []VPN{
						{VPNID: vpnID, Name: "vpn1", EnableHTTPS: true, CertificateID: nil},
					}

					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).AnyTimes()
					mocks.vpn.EXPECT().
						Get(t.Context(), vpnID).
						Return(&vpn.VPN{Enabled: true, Driver: "driver1"}, nil)
					mocks.vpn.EXPECT().
						GetAvailableDrivers(t.Context()).
						Return([]vpn.AvailableDriver{
							{
								ID:                 "driver1",
								EndpointSSLSupport: vpn.DriverManagedEndpointSSLSupport,
							},
						}, nil)

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostVpnCertificateRequired)
				})

				t.Run("driver managed - certificate not found", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					vpnID := uuid.New()
					certID := uuid.New()
					h.VPNs = []VPN{
						{VPNID: vpnID, Name: "vpn1", EnableHTTPS: true, CertificateID: &certID},
					}

					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).AnyTimes()
					mocks.vpn.EXPECT().
						Get(t.Context(), vpnID).
						Return(&vpn.VPN{Enabled: true, Driver: "driver1"}, nil)
					mocks.vpn.EXPECT().
						GetAvailableDrivers(t.Context()).
						Return([]vpn.AvailableDriver{
							{
								ID:                 "driver1",
								EndpointSSLSupport: vpn.DriverManagedEndpointSSLSupport,
							},
						}, nil)
					mocks.certificate.EXPECT().Exists(t.Context(), certID).Return(false, nil)

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostVpnCertificateNotFound)
				})

				t.Run("provider managed - certificate prohibited", func(t *testing.T) {
					hostValidator, mocks := setupValidator(t)
					h := newHost()
					vpnID := uuid.New()
					certID := uuid.New()
					h.VPNs = []VPN{
						{VPNID: vpnID, Name: "vpn1", EnableHTTPS: true, CertificateID: &certID},
					}

					mocks.binding.EXPECT().
						Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
						Return(nil).AnyTimes()
					mocks.vpn.EXPECT().
						Get(t.Context(), vpnID).
						Return(&vpn.VPN{Enabled: true, Driver: "driver1"}, nil)
					mocks.vpn.EXPECT().
						GetAvailableDrivers(t.Context()).
						Return([]vpn.AvailableDriver{
							{
								ID:                 "driver1",
								EndpointSSLSupport: vpn.ProviderManagedEndpointSSLSupport,
							},
						}, nil)

					err := hostValidator.validate(t.Context(), h)
					assertViolations(t, err, i18n.K.CoreHostVpnCertificateProhibited)
				})
			})
		})

		t.Run("validates ACLs", func(t *testing.T) {
			hostValidator, mocks := setupValidator(t)
			h := newHost()
			aclID := uuid.New()
			h.AccessListID = &aclID

			mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
			mocks.binding.EXPECT().
				Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			mocks.accessList.EXPECT().Exists(t.Context(), aclID).Return(false, nil)

			err := hostValidator.validate(t.Context(), h)
			assertViolations(t, err, i18n.K.CoreHostAccessListNotFound)
		})

		t.Run("validates Cache", func(t *testing.T) {
			hostValidator, mocks := setupValidator(t)
			h := newHost()
			cacheID := uuid.New()
			h.CacheID = &cacheID

			mocks.vpn.EXPECT().GetAvailableDrivers(t.Context()).Return(nil, nil).AnyTimes()
			mocks.binding.EXPECT().
				Validate(t.Context(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
				Return(nil)
			mocks.cache.EXPECT().Exists(t.Context(), cacheID).Return(false, nil)

			err := hostValidator.validate(t.Context(), h)
			assertViolations(t, err, i18n.K.CoreHostCacheNotFound)
		})
	})
}
