package utils

import (
	"fmt"
	"strings"
)

// Add suffix if it doesn't exist
//
// For example,
// "str"  -> "str_";
// "str_" -> "str_";
func AddSuffix(s string, suffix string) string {
	if !strings.HasSuffix(s, suffix) {
		return fmt.Sprintf("%s%s", s, suffix)
	}
	return s
}
