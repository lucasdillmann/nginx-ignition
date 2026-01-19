package validator

import (
	"fmt"
	"strings"

	"dillmann.com.br/nginx-ignition/core/common/log"
	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

var (
	duplicatedValuesBypassPrefixes = []string{
		"certificate/",
	}
)

func Validate(files []reader.PropertiesFile) []string {
	keyUsage := make(map[string]map[string]bool)
	valueUsage := make(map[string]map[string][]string)
	keyCounts := make(map[string]map[string]int)
	allLanguageTags := make([]string, 0, len(files))

	for _, file := range files {
		allLanguageTags = append(allLanguageTags, file.LanguageTag)

		for _, msg := range file.Messages {
			trackKeyUsage(keyUsage, msg.PropertiesKey, file.LanguageTag)
			trackValueUsage(valueUsage, msg.Value, file.LanguageTag, msg.PropertiesKey)
			trackKeyCount(keyCounts, msg.PropertiesKey, file.LanguageTag)
		}
	}

	var violations []string
	violations = append(violations, checkDuplicateKeys(keyCounts)...)
	violations = append(violations, checkMissingKeys(keyUsage, allLanguageTags, len(files))...)
	violations = append(violations, checkDuplicateValues(valueUsage)...)

	return violations
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

func trackKeyCount(counts map[string]map[string]int, key, tag string) {
	if counts[tag] == nil {
		counts[tag] = make(map[string]int)
	}

	counts[tag][key]++
}

func checkDuplicateKeys(keyCounts map[string]map[string]int) []string {
	var violations []string

	for tag, keys := range keyCounts {
		for key, count := range keys {
			if count > 1 {
				violations = append(violations, fmt.Sprintf("Key '%s' is duplicated in language '%s'", key, tag))
			}
		}
	}

	return violations
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

func checkDuplicateValues(valueUsage map[string]map[string][]string) []string {
	output := make([]string, 0)

	for val, languages := range valueUsage {
		for tag, keys := range languages {
			if len(keys) == 1 {
				continue
			}

			if bypassDuplicatedValueCheck(keys[0]) {
				log.Warnf("Duplicate value '%s' ignored: key bypassed", val)
				continue
			}

			output = append(
				output,
				fmt.Sprintf("Value '%s' is duplicated in language '%s' on keys %v", val, tag, keys),
			)
		}
	}

	return output
}

func bypassDuplicatedValueCheck(key string) bool {
	for _, prefix := range duplicatedValuesBypassPrefixes {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}

	return false
}
