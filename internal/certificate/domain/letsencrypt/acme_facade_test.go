package letsencrypt

import (
	"context"
	"errors"
	"testing"

	"github.com/go-acme/lego/v5/challenge/dns01"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/lucasdillmann/nginx-ignition/internal/business/domain/certificate"
)

func Test_isValidDNSServer(t *testing.T) {
	tests := []struct {
		value    string
		expected bool
	}{
		{value: "1.1.1.1", expected: true},
		{value: "8.8.8.8:53", expected: true},
		{value: "2001:4860:4860::8888", expected: true},
		{value: "[2001:4860:4860::8888]:53", expected: true},
		{value: "", expected: false},
		{value: "example.com", expected: false},
		{value: "999.1.1.1", expected: false},
		{value: "1.1.1.1:0", expected: false},
		{value: "1.1.1.1:65536", expected: false},
		{value: "1.1.1.1:abc", expected: false},
	}

	for _, test := range tests {
		t.Run(test.value, func(t *testing.T) {
			assert.Equal(t, test.expected, isValidDNSServer(test.value))
		})
	}
}

func Test_configureCustomDNSServers(t *testing.T) {
	tests := []struct {
		parameters   map[string]any
		name         string
		expectError  bool
		expectChange bool
	}{
		{
			name: "enabled",
			parameters: map[string]any{
				"overrideDnsServers": true,
				"customDnsServers":   "1.1.1.1\n8.8.8.8:53",
			},
			expectChange: true,
		},
		{
			name: "disabled",
			parameters: map[string]any{
				"overrideDnsServers": false,
				"customDnsServers":   "1.1.1.1",
			},
		},
		{
			name: "propagation checks skipped",
			parameters: map[string]any{
				"bypassDnsPropagationChecks": true,
				"overrideDnsServers":         true,
				"customDnsServers":           "1.1.1.1",
			},
		},
		{
			name: "invalid address",
			parameters: map[string]any{
				"overrideDnsServers": true,
				"customDnsServers":   "example.com",
			},
			expectError: true,
		},
		{
			name: "empty list",
			parameters: map[string]any{
				"overrideDnsServers": true,
				"customDnsServers":   "",
			},
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			previous := dns01.DefaultClient()
			t.Cleanup(func() { dns01.SetDefaultClient(previous) })

			err := configureCustomDNSServers(t.Context(), test.parameters)

			if test.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if test.expectChange {
				assert.NotSame(t, previous, dns01.DefaultClient())
			} else {
				assert.Same(t, previous, dns01.DefaultClient())
			}
		})
	}
}

func Test_withCustomDNSServers_restoresPreviousClient(t *testing.T) {
	previous := dns01.DefaultClient()
	t.Cleanup(func() { dns01.SetDefaultClient(previous) })

	expected := &certificate.Certificate{}

	result, err := withCustomDNSServers(
		t.Context(),
		customDNSServersParameters(),
		func() (*certificate.Certificate, error) {
			assert.NotSame(t, previous, dns01.DefaultClient())

			return expected, nil
		},
	)

	require.NoError(t, err)
	assert.Same(t, expected, result)
	assert.Same(t, previous, dns01.DefaultClient())
}

func Test_withCustomDNSServers_restoresPreviousClientOnError(t *testing.T) {
	previous := dns01.DefaultClient()
	t.Cleanup(func() { dns01.SetDefaultClient(previous) })

	expectedErr := errors.New("issuance failed")

	result, err := withCustomDNSServers(
		t.Context(),
		customDNSServersParameters(),
		func() (*certificate.Certificate, error) {
			return nil, expectedErr
		},
	)

	require.ErrorIs(t, err, expectedErr)
	assert.Nil(t, result)
	assert.Same(t, previous, dns01.DefaultClient())
}

func Test_acquireIssuanceLock_respectsCanceledContext(t *testing.T) {
	require.NoError(t, acquireIssuanceLock(t.Context()))
	defer releaseIssuanceLock()

	ctx, cancel := context.WithCancel(t.Context())
	cancel()

	assert.ErrorIs(t, acquireIssuanceLock(ctx), context.Canceled)
}

func customDNSServersParameters() map[string]any {
	return map[string]any{
		"overrideDnsServers": true,
		"customDnsServers":   "1.1.1.1",
	}
}
