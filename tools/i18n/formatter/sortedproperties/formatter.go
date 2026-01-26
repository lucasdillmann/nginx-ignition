package sortedproperties

import (
	"fmt"
	"io"
	"slices"

	"dillmann.com.br/nginx-ignition/tools/i18n/formatter"
	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

type sortedFormatter struct{}

func New() formatter.Formatter {
	return &sortedFormatter{}
}

func (*sortedFormatter) Format(propertiesFile *reader.PropertiesFile, writer io.Writer) error {
	lines := make([]string, 0, len(propertiesFile.Messages))
	for _, msg := range propertiesFile.Messages {

		lines = append(lines, fmt.Sprintf("%s=%s", msg.PropertiesKey, msg.RawValue))
	}

	slices.Sort(lines)

	for _, line := range lines {
		_, err := writer.Write([]byte(line + "\n"))
		if err != nil {
			return err
		}
	}

	return nil
}
