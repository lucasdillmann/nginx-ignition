package validator

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

const (
	duplicatedValuesCheckEnabled = false
)

func Validate(files []reader.PropertiesFile) []string {
	keyUsage := make(map[string]map[string]bool)
	valueUsage := make(map[string]map[string][]string)
	allLanguageTags := make([]string, 0, len(files))

	for _, file := range files {
		allLanguageTags = append(allLanguageTags, file.LanguageTag)

		for _, msg := range file.Messages {
			trackKeyUsage(keyUsage, msg.PropertiesKey, file.LanguageTag)
			trackValueUsage(valueUsage, msg.Value, file.LanguageTag, msg.PropertiesKey)
		}
	}

	checkDuplicateValues(valueUsage)

	return checkMissingKeys(keyUsage, allLanguageTags, len(files))
}

func trackKeyUsage(usage map[string]map[string]bool, key, tag string) {
	if usage[key] == nil {
		usage[key] = make(map[string]bool)
	}

	usage[key][tag] = true
}

func trackValueUsage(usage map[string]map[string][]string, val, tag, key string) {
	if usage[val] == nil {
		usage[val] = make(map[string][]string)
	}

	usage[val][tag] = append(usage[val][tag], key)
}

func checkMissingKeys(keyUsage map[string]map[string]bool, tags []string, expectedCount int) []string {
	var violations []string

	for key, languagesWithKey := range keyUsage {
		if len(languagesWithKey) >= expectedCount {
			continue
		}

		for _, tag := range tags {
			if !languagesWithKey[tag] {
				violations = append(violations, fmt.Sprintf("Key '%s' is missing in language '%s'", key, tag))
			}
		}
	}

	return violations
}

func checkDuplicateValues(valueUsage map[string]map[string][]string) {
	for val, languages := range valueUsage {
		for tag, keys := range languages {
			if len(keys) > 1 {
				log.Warnf("Language %s has duplicate value '%s' in keys %v", tag, val, keys)
			}
		}
	}
}
