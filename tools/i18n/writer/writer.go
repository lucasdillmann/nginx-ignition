package writer

import (
	"fmt"
	"os"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter"
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter/golangdictionary"
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter/golangkeys"
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter/sortedproperties"
	"dillmann.com.br/nginx-ignition/tools/i18n/formatter/typescriptkeys"
	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

func Write(propertiesFiles []reader.PropertiesFile) error {
	if err := writeFile(
		propertiesFiles[0],
		typescriptkeys.New(),
		"frontend/src/core/i18n/model/MessageKey.generated.ts",
	); err != nil {
		return err
	}

	if err := writeFile(
		propertiesFiles[0],
		golangkeys.New(),
		"i18n/keys.generated.go",
	); err != nil {
		return err
	}

	for _, propertiesFile := range propertiesFiles {
		if err := writeFile(
			propertiesFile,
			golangdictionary.New(),
			fmt.Sprintf("i18n/%s.generated.go", strings.ToLower(propertiesFile.NormalizedLanguageTag)),
		); err != nil {
			return err
		}

		if err := writeFile(
			propertiesFile,
			sortedproperties.New(),
			fmt.Sprintf("i18n/messages-%s.properties", strings.ToLower(propertiesFile.LanguageTag)),
		); err != nil {
			return err
		}
	}

	return nil
}

func writeFile(
	propertiesFile reader.PropertiesFile,
	fmt formatter.Formatter,
	path string,
) error {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	if err := fmt.Format(&propertiesFile, file); err != nil {
		return err
	}

	log.Infof("i18n code file written: %s", path)
	return nil
}
