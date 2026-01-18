package validator

import (
	"fmt"

	"dillmann.com.br/nginx-ignition/tools/i18n/reader"
)

func Validate(files []reader.PropertiesFile) []string {
	var violations []string

	keyUsage := make(map[string]map[string]bool)
	valueUsage := make(map[string]map[string][]string)
	allLanguageTags := make([]string, 0, len(files))

	for _, file := range files {
		allLanguageTags = append(allLanguageTags, file.LanguageTag)

		for _, msg := range file.Messages {
			if keyUsage[msg.PropertiesKey] == nil {
				keyUsage[msg.PropertiesKey] = make(map[string]bool)
			}
			keyUsage[msg.PropertiesKey][file.LanguageTag] = true

			if valueUsage[msg.Value] == nil {
				valueUsage[msg.Value] = make(map[string][]string)
			}
			valueUsage[msg.Value][file.LanguageTag] = append(valueUsage[msg.Value][file.LanguageTag], msg.PropertiesKey)
		}
	}

	for key, languagesWithKey := range keyUsage {
		if len(languagesWithKey) < len(files) {
			for _, tag := range allLanguageTags {
				if !languagesWithKey[tag] {
					violation := fmt.Sprintf("Key '%s' is missing in language '%s'", key, tag)
					violations = append(violations, violation)
				}
			}
		}
	}

	//for val, languages := range valueUsage {
	//	for tag, keys := range languages {
	//		if len(keys) > 1 {
	//			violation := fmt.Sprintf("Language '%s' has duplicate value '%s' in keys: %v", tag, val, keys)
	//			violations = append(violations, violation)
	//		}
	//	}
	//}

	return violations
}
