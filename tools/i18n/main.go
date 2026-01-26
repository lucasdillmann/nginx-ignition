package main

import (
	"os"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
	"dillmann.com.br/nginx-ignition/tools/i18n/validator"
	"dillmann.com.br/nginx-ignition/tools/i18n/writer"
)

func main() {
	log.EnableStackTrace(false)
	log.Info("Starting i18n code generation...")

	propertiesFiles, err := reader.ReadPropertiesFiles("i18n")
	if err != nil {
		log.Errorf("Error reading properties files: %v", err)
		os.Exit(1)
	}

	if len(propertiesFiles) == 0 {
		log.Error("No properties files found in the i18n folder")
		os.Exit(1)
	}

	problems := validator.Validate(propertiesFiles)
	if len(problems) > 0 {
		problemsMerged := strings.Join(problems, "\n- ")
		log.Errorf("One or more problems were found in the properties files: \n- %s", problemsMerged)
		os.Exit(1)
	}

	if err := writer.Write(propertiesFiles); err != nil {
		log.Errorf("Error writing files: %v", err)
		os.Exit(1)
	}

	log.Info("Code generation for the i18n messages completed successfully")
}
