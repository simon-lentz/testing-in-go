package underscore

import (
	"strings"
)

// If no test file has been made for the underscore.go file,
// highlighting the function Camel(), then selecting
// "Go: Generate Units Tests For Function" will generate
// a new underscore_test.go file with the table-test boilerplate.
func Camel(str string) string {
	var sb strings.Builder
	for _, ch := range str {
		// If the character is a capital letter, prepend an underscore.
		if ch >= 'A' && ch <= 'Z' {
			sb.WriteRune('_')
		}
		sb.WriteRune(ch)
	}
	return strings.ToLower(sb.String())
}
