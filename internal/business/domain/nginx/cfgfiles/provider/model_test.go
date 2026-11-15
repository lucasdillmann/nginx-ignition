package provider

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_File(t *testing.T) {
	t.Run("FormattedContents", func(t *testing.T) {
		t.Run("indents content correctly", func(t *testing.T) {
			file := File{
				Contents: `
				http {
					server {
						listen 80;
					}
				}`,
			}

			expected := "http {\n    server {\n        listen 80;\n    }\n}"
			assert.Equal(t, expected, file.FormattedContents())
		})

		t.Run("handles multiple levels of indentation", func(t *testing.T) {
			file := File{
				Contents: `
				stream {
					upstream backend {
						server 127.0.0.1:8080;
					}
				}`,
			}

			expected := "stream {\n    upstream backend {\n        server 127.0.0.1:8080;\n    }\n}"
			assert.Equal(t, expected, file.FormattedContents())
		})

		t.Run("handles empty contents", func(t *testing.T) {
			file := File{Contents: ""}
			assert.Equal(t, "", file.FormattedContents())
		})

		t.Run("handles single line without braces", func(t *testing.T) {
			file := File{Contents: "listen 80;"}
			assert.Equal(t, "listen 80;", file.FormattedContents())
		})

		t.Run("handles braces on same line", func(t *testing.T) {
			file := File{Contents: "location / { return 200; }"}
			expected := "location / { return 200; }"
			assert.Equal(t, expected, file.FormattedContents())
		})

		t.Run("handles unbalanced closing braces", func(t *testing.T) {
			file := File{Contents: "}\n}"}
			assert.Equal(t, "}\n}", file.FormattedContents())
		})

		t.Run("handles comments", func(t *testing.T) {
			file := File{
				Contents: `
				# This is a comment
				http {
					# Inner comment
					listen 80;
				}`,
			}
			expected := "# This is a comment\nhttp {\n    # Inner comment\n    listen 80;\n}"
			assert.Equal(t, expected, file.FormattedContents())
		})
	})

	t.Run("Bytes", func(t *testing.T) {
		t.Run("formats textual contents", func(t *testing.T) {
			file := File{Contents: "\nhttp {\nserver {\nlisten 80;\n}\n}"}
			assert.Equal(
				t,
				[]byte("http {\n    server {\n        listen 80;\n    }\n}"),
				file.Bytes(),
			)
		})

		t.Run("returns binary contents verbatim", func(t *testing.T) {
			binary := "\x00\x01\x0a\x20\x0b\x0c\x0d\x0e"
			file := File{Contents: binary, Binary: true}
			assert.Equal(t, []byte(binary), file.Bytes())
		})
	})
}
