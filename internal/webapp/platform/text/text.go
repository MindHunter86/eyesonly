package text

import (
	"strings"
	"unicode/utf8"
)

func Preview(value string, limit int) string {
	value = strings.Join(strings.Fields(value), " ")
	if value == "" {
		return "empty secret"
	}
	if limit <= 0 {
		limit = 64
	}
	if utf8.RuneCountInString(value) <= limit {
		return value
	}
	runes := []rune(value)
	return string(runes[:limit]) + "..."
}
