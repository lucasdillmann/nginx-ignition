package formatter

import (
	"io"

	"nginx-ignition/tools/i18n/reader"
)

type Formatter interface {
	Format(*reader.PropertiesFile, io.Writer) error
}
