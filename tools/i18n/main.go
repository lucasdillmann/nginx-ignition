package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"nginx-ignition/tools/i18n/reader"
	"nginx-ignition/tools/i18n/validator"
	"nginx-ignition/tools/i18n/writer"
)

func main() {
	log.Println("Starting i18n code generation...")

	if _, err := os.Stat("internal/i18n"); os.IsNotExist(err) {
		if err := os.Chdir(filepath.Join("..", "..")); err != nil {
			log.Fatal(err)
		}
	}

	propertiesFiles, err := reader.ReadPropertiesFiles("internal/i18n")
	if err != nil {
		log.Printf("Error reading properties files: %v", err)
		os.Exit(1)
	}

	if len(propertiesFiles) == 0 {
		log.Println("No properties files found in the i18n folder")
		os.Exit(1)
	}

	problems := validator.Validate(propertiesFiles)
	if len(problems) > 0 {
		problemsMerged := strings.Join(problems, "\n- ")
		log.Printf("One or more problems were found in the properties files: \n- %s", problemsMerged)
		os.Exit(1)
	}

	if err := writer.Write(propertiesFiles); err != nil {
		log.Printf("Error writing files: %v", err)
		os.Exit(1)
	}

	log.Println("Code generation for the i18n messages completed successfully")
}
