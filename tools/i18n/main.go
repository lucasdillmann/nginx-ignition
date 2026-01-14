package main

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

const (
	baseDir = "i18n"
)

func main() {
	log.Infof("Starting i18n code generation...")
	files, err := os.ReadDir(baseDir)
	if err != nil {
		panic(err)
	}

	allKeys := make(map[string]bool)
	propertiesByLang := make(map[string]properties)

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".properties" {
			continue
		}

		log.Infof("i18n file found: %s", file.Name())

		lang := strings.TrimPrefix(strings.TrimSuffix(file.Name(), ".properties"), "messages-")
		propertiesByLang[lang] = readPropertiesFile(filepath.Join(baseDir, file.Name()), &allKeys)
	}

	validateKeys(allKeys, propertiesByLang)

	log.Infof("Writing keys file...")
	writeKeysFile(baseDir, allKeys)

	for lang, props := range propertiesByLang {
		log.Infof("Writing %s dictionary file...", lang)
		writeDictionaryFile(baseDir, lang, props)
	}

	log.Infof("i18n code generation completed")
}

func toPascalCase(s string, caser cases.Caser) string {
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '.' || r == '-' || r == '_'
	})
	for index := range parts {
		parts[index] = caser.String(parts[index])
	}

	return strings.Join(parts, "")
}
