package cfgfiles

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFile_FormattedContents(t *testing.T) {
	t.Run("indents content correctly", func(t *testing.T) {
		f := File{
			Contents: `
			http {
				server {
					listen 80;
				}
			}`,
		}

		expected := "http {\n    server {\n        listen 80;\n    }\n}"
		assert.Equal(t, expected, f.FormattedContents())
	})

	t.Run("handles multiple levels of indentation", func(t *testing.T) {
		f := File{
			Contents: `
			stream {
				upstream backend {
					server 127.0.0.1:8080;
				}
			}`,
		}

		expected := "stream {\n    upstream backend {\n        server 127.0.0.1:8080;\n    }\n}"
		assert.Equal(t, expected, f.FormattedContents())
	})
}
