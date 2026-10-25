package letsencrypt

import (
	"testing"

	"github.com/stretchr/testify/assert"
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
