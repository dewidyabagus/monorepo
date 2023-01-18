package mstring

import "strings"

func ContainsOneOf(s string, values ...string) bool {
	if len(values) == 0 {
		return false
	}

	for _, substr := range values {
		if strings.Contains(s, substr) {
			return true
		}
	}

	return false
}
