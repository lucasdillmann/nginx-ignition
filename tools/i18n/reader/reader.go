package reader

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"dillmann.com.br/nginx-ignition/core/common/log"
)

func ReadPropertiesFiles(baseFolder string) ([]PropertiesFile, error) {
	folderContents, err := os.ReadDir(baseFolder)
	if err != nil {
		return nil, err
	}

	output := make([]PropertiesFile, 0)

	for _, fileInfo := range folderContents {
		if fileInfo.IsDir() ||
			!strings.HasSuffix(fileInfo.Name(), ".properties") ||
			!strings.HasPrefix(fileInfo.Name(), "messages_") {
			continue
		}

		fullPath := fmt.Sprintf("%s%c%s", baseFolder, filepath.Separator, fileInfo.Name())
		file, err := os.Open(fullPath)
		if err != nil {
			return nil, err
		}

		defer file.Close()

		propertiesFile, err := readPropertiesFile(file)
		if err != nil {
			return nil, err
		}

		log.Infof("i18n properties file foud for %s: %s", propertiesFile.LanguageTag, fullPath)
		output = append(output, propertiesFile)
	}

	return output, nil
}

func readPropertiesFile(file *os.File) (PropertiesFile, error) {
	baseName := filepath.Base(file.Name())
	languageTag := strings.TrimSuffix(baseName, ".properties")
	languageTag = strings.TrimPrefix(languageTag, "messages_")

	scanner := bufio.NewScanner(file)
	properties := PropertiesFile{
		LanguageTag:           languageTag,
		NormalizedLanguageTag: normalizeLanguageTag(languageTag),
		Messages:              make([]Message, 0),
	}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") {
			continue
		}

		if strings.TrimSpace(line) == "" {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			propertiesKey := strings.TrimSpace(parts[0])
			properties.Messages = append(properties.Messages, Message{
				PropertiesKey: propertiesKey,
				CamelCaseKey:  camelCaseKey(propertiesKey),
				SnakeCaseKey:  snakeCaseKey(propertiesKey),
				RawValue:      parts[1],
				Value:         normalizeValue(parts[1]),
			})
		}
	}

	return properties, scanner.Err()
}

func normalizeValue(value string) string {
	output := strings.TrimSpace(value)
	output = strings.ReplaceAll(value, "\n", "\\n")
	output = strings.ReplaceAll(output, "\r", "")
	return strings.ReplaceAll(output, "\"", "\\\"")
}

func normalizeLanguageTag(languageTag string) string {
	parts := strings.Split(languageTag, "-")
	firstPart := cases.Title(language.English).String(parts[0])

	if len(parts) == 1 {
		return firstPart
	}

	secondPart := strings.ToUpper(parts[1])
	return fmt.Sprintf("%s_%s", firstPart, secondPart)
}

func camelCaseKey(key string) string {
	parts := splitKey(key)
	caser := cases.Title(language.English)

	var result string
	for _, part := range parts {
		result += caser.String(part)
	}

	return result
}

func snakeCaseKey(key string) string {
	parts := splitKey(key)
	for index := range parts {
		parts[index] = strings.ToLower(parts[index])
	}

	return strings.Join(parts, "_")
}

func splitKey(key string) []string {
	f := func(r rune) bool {
		return r == '/' || r == '-' || r == '_'
	}

	return strings.FieldsFunc(key, f)
}
