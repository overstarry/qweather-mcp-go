package utils

import (
	"strings"
)

// JoinStrings Join string array and filter empty values
func JoinStrings(strs []string, sep string) string {
	var filtered []string
	for _, s := range strs {
		if s != "" {
			filtered = append(filtered, s)
		}
	}
	return strings.Join(filtered, sep)
}
