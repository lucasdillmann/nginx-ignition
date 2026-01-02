package binding

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"dillmann.com.br/nginx-ignition/core/certificate"
)

func Test_BuildCommands(t *testing.T) {
	t.Run("builds commands with all service methods", func(t *testing.T) {
		certCommands := &certificate.Commands{}
		commands := buildCommands(certCommands)

		assert.NotNil(t, commands)
		assert.NotNil(t, commands.Validate)
	})
}
