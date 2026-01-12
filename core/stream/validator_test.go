package stream

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"dillmann.com.br/nginx-ignition/core/common/i18n"
	"dillmann.com.br/nginx-ignition/core/common/ptr"
	"dillmann.com.br/nginx-ignition/core/common/validation"
)

func Test_validator(t *testing.T) {
	validate := func(s *Stream) error {
		return newValidator().validate(t.Context(), s)
	}

	assertViolations := func(t *testing.T, err error, msgs ...string) {
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

	t.Run("validate", func(t *testing.T) {
		t.Run("returns error when stream is nil", func(t *testing.T) {
			err := validate(nil)
			assertViolations(t, err, i18n.K.StreamValidationNilStream)
		})

		t.Run("valid stream passes", func(t *testing.T) {
			s := newStream()
			require.NoError(t, validate(s))
		})

		t.Run("validates name", func(t *testing.T) {
			s := newStream()
			s.Name = ""
			err := validate(s)
			assertViolations(t, err, i18n.K.CommonValidationCannotBeEmpty)
		})

		t.Run("validates type", func(t *testing.T) {
			s := newStream()
			s.Type = "INVALID"
			err := validate(s)
			assertViolations(t, err, i18n.K.CommonValidationInvalidValue)
		})

		t.Run("validates binding", func(t *testing.T) {
			t.Run("address required", func(t *testing.T) {
				s := newStream()
				s.Binding.Address = ""
				err := validate(s)
				assertViolations(t, err, i18n.K.CommonValidationCannotBeEmpty)
			})

			t.Run("protocol required", func(t *testing.T) {
				s := newStream()
				s.Binding.Protocol = "INVALID"
				err := validate(s)
				assertViolations(t, err, i18n.K.CommonValidationInvalidValue)
			})

			t.Run("port required for TCP/UDP", func(t *testing.T) {
				s := newStream()
				s.Binding.Protocol = TCPProtocol
				s.Binding.Port = nil
				err := validate(s)
				assertViolations(t, err, i18n.K.StreamValidationPortRequired)
			})

			t.Run("port range validation", func(t *testing.T) {
				s := newStream()
				s.Binding.Protocol = TCPProtocol
				s.Binding.Port = ptr.Of(70000)
				err := validate(s)
				assertViolations(t, err, i18n.K.CommonValidationBetweenValues)
			})

			t.Run("socket protocol validation", func(t *testing.T) {
				s := newStream()
				s.Binding.Protocol = SocketProtocol
				s.Binding.Address = "invalid" // Missing /
				s.Binding.Port = ptr.Of(80)   // Should be nil
				err := validate(s)
				assertViolations(
					t,
					err,
					i18n.K.CommonValidationStartsWithSlashRequired,
					i18n.K.StreamValidationPortNotAllowedForSocket,
				)
			})
		})

		t.Run("validates default backend", func(t *testing.T) {
			t.Run("validates address", func(t *testing.T) {
				s := newStream()
				s.DefaultBackend.Address.Address = ""
				err := validate(s)
				assertViolations(t, err, i18n.K.CommonValidationCannotBeEmpty)
			})

			t.Run("validates circuit breaker", func(t *testing.T) {
				s := newStream()
				s.DefaultBackend.CircuitBreaker = &CircuitBreaker{
					MaxFailures: 0,
					OpenSeconds: -1,
				}
				err := validate(s)
				assertViolations(
					t,
					err,
					i18n.K.CommonValidationCannotBeZero,
					i18n.K.CommonValidationCannotBeZero,
				)
			})

			t.Run("nil circuit breaker is valid", func(t *testing.T) {
				s := newStream()
				s.DefaultBackend.CircuitBreaker = nil
				require.NoError(t, validate(s))
			})
		})

		t.Run("validates routes", func(t *testing.T) {
			t.Run("skipped for SimpleType", func(t *testing.T) {
				s := newStream()
				s.Type = SimpleType
				s.Routes = nil
				require.NoError(t, validate(s))
			})

			t.Run("required for SNIRouterType", func(t *testing.T) {
				s := newStream()
				s.Type = SNIRouterType
				s.Routes = nil
				err := validate(s)
				assertViolations(t, err, i18n.K.StreamValidationRoutesRequiredForSni)
			})

			t.Run("validates route entries", func(t *testing.T) {
				s := newStream()
				s.Type = SNIRouterType
				s.Routes = []Route{
					{
						DomainNames: []string{},
						Backends:    []Backend{},
					},
				}
				err := validate(s)
				assertViolations(
					t,
					err,
					i18n.K.StreamValidationAtLeastOneDomain,
					i18n.K.StreamValidationAtLeastOneBackend,
				)
			})

			t.Run("validates domain names", func(t *testing.T) {
				s := newStream()
				s.Type = SNIRouterType
				s.Routes = []Route{
					{
						DomainNames: []string{""},
						Backends: []Backend{
							{
								Address: Address{
									Protocol: TCPProtocol,
									Address:  "127.0.0.1",
									Port:     ptr.Of(80),
								},
							},
						},
					},
				}
				err := validate(s)
				assertViolations(t, err, i18n.K.CommonValidationCannotBeEmpty)

				s.Routes[0].DomainNames[0] = "invalid_domain"
				err = validate(s)
				assertViolations(t, err, i18n.K.CommonValidationInvalidDomainName)
			})

			t.Run("validates backends", func(t *testing.T) {
				s := newStream()
				s.Type = SNIRouterType
				s.Routes = []Route{
					{
						DomainNames: []string{"example.com"},
						Backends: []Backend{
							{Address: Address{Address: ""}},
						},
					},
				}
				err := validate(s)
				assertViolations(t, err, i18n.K.CommonValidationCannotBeEmpty)
			})
		})

		t.Run("validates feature set", func(t *testing.T) {
			t.Run("allows TCP features for TCP protocol", func(t *testing.T) {
				s := newStream()
				s.Binding.Protocol = TCPProtocol
				s.FeatureSet = FeatureSet{
					TCPKeepAlive: true,
					TCPNoDelay:   true,
					TCPDeferred:  true,
				}
				require.NoError(t, validate(s))
			})

			t.Run("disallows TCP features for non-TCP protocol", func(t *testing.T) {
				s := newStream()
				s.Binding.Protocol = UDPProtocol
				s.FeatureSet = FeatureSet{
					TCPKeepAlive: true,
					TCPNoDelay:   true,
					TCPDeferred:  true,
				}
				err := validate(s)
				assertViolations(t, err,
					i18n.K.StreamValidationFeatureOnlyForTCP,
					i18n.K.StreamValidationFeatureOnlyForTCP,
					i18n.K.StreamValidationFeatureOnlyForTCP,
				)
			})
		})
	})

	t.Run("validateName", func(t *testing.T) {
		streamValidator := newValidator()
		s := newStream()

		s.Name = strings.Repeat("a", 256)
		streamValidator.validateName(t.Context(), s)
		assert.Empty(t, streamValidator.delegate.Result())

		s.Name = "   "
		streamValidator.validateName(t.Context(), s)
		assert.Error(t, streamValidator.delegate.Result())
	})
}
