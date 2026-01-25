package validator

import (
	"fmt"
	"path"
	"regexp"
	"slices"
	"sort"
	"strings"

	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

var (
	duplicatedValuesBypassPrefixes = []string{
		"certificate/",
		"frontend/components/domainnames/placeholder",
		"frontend/stream/components/routeform/domains",
		"frontend/settings/tabs/nginx/logs/server",
		"frontend/user/menu/change-password-title",
		"frontend/accesslist/section-credentials",
		"frontend/nginx/control/feature-streams",
		"frontend/settings/tabs/nginx/log-level-alert",
		"frontend/settings/tabs/nginx/log-level-warn",
		"frontend/stream/components/routeform/backends",
		"frontend/stream/form/backend-title",
		"frontend/logs/host-logs",
		"frontend/logs/server-logs",
		"frontend/host/components/hostbindings/protocol",
		"frontend/settings/tabs/ignition/time-unit-minutes",
	}
	placeholderRegex = regexp.MustCompile(`\$\{[a-zA-Z0-9_\-]+\}`)
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
	violations = append(violations, checkBlankStrings(files)...)
	violations = append(violations, checkDuplicateKeys(keyCounts)...)
	violations = append(violations, checkMissingKeys(keyUsage, allLanguageTags, len(files))...)
	violations = append(violations, checkDuplicateValues(valueUsage)...)
	violations = append(violations, checkPlaceholders(files)...)

	return violations
}

func checkBlankStrings(files []reader.PropertiesFile) []string {
	var violations []string

	for _, file := range files {
		for _, msg := range file.Messages {
			if strings.TrimSpace(msg.PropertiesKey) == "" {
				violations = append(
					violations,
					fmt.Sprintf("Blank key found in language '%s'", file.LanguageTag),
				)
			}

			if strings.TrimSpace(msg.Value) == "" {
				violations = append(
					violations,
					fmt.Sprintf("Blank value found for key '%s' in language '%s'", msg.PropertiesKey, file.LanguageTag),
				)
			}
		}
	}

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

			if bypassDuplicatedValueCheck(keys) {
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

func bypassDuplicatedValueCheck(keys []string) bool {
	for _, key := range keys {
		for _, prefix := range duplicatedValuesBypassPrefixes {
			if strings.HasPrefix(key, prefix) {
				return true
			}
		}
	}

	return areKeysSingularPluralVariants(keys)
}

func areKeysSingularPluralVariants(keys []string) bool {
	normalized := make(map[string]bool)

	for _, key := range keys {
		base := path.Base(key)
		stem := strings.TrimSuffix(base, "s")
		normalized[stem] = true
	}

	return len(normalized) == 1
}

func checkPlaceholders(files []reader.PropertiesFile) []string {
	var violations []string
	type placeholderInfo struct {
		placeholders []string
		lang         string
	}
	refs := make(map[string]placeholderInfo)

	for _, file := range files {
		for _, msg := range file.Messages {
			matches := placeholderRegex.FindAllString(msg.Value, -1)
			sort.Strings(matches)

			if ref, exists := refs[msg.PropertiesKey]; exists {
				if !slices.Equal(ref.placeholders, matches) {
					violations = append(
						violations,
						fmt.Sprintf(
							"Placeholder mismatch for key '%s': Language '%s' has %v, but Language '%s' has %v",
							msg.PropertiesKey, file.LanguageTag, matches, ref.lang, ref.placeholders,
						),
					)
				}
			} else {
				refs[msg.PropertiesKey] = placeholderInfo{matches, file.LanguageTag}
			}
		}
	}

	return violations
}
