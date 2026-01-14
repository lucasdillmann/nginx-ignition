package main

import (
	"fmt"
)

func validateKeys(allKeys map[string]bool, langProps map[string]properties) {
	for lang, props := range langProps {
		for key := range allKeys {
			if _, exists := props[key]; !exists {
				panic(fmt.Sprintf("Missing key %s in %s", key, lang))
			}
		}
	}
}
