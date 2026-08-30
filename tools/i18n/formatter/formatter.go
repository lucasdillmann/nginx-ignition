package formatter

import (
	"io"

	"github.com/lucasdillmann/nginx-ignition/tools/i18n/reader"
)

type Formatter interface {
	Format(*reader.PropertiesFile, io.Writer) error
}
