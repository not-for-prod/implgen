package strtools

import "strings"

// SnakeCase a helper function to convert camelCase to snake_case
func SnakeCase(str string) string {
	var result []rune
	for i, c := range str {
		if i > 0 && 'A' <= c && c <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, rune(strings.ToLower(string(c))[0]))
	}
	return string(result)
}

func KebabCase(str string) string {
	var result []rune
	for i, c := range str {
		if i > 0 && 'A' <= c && c <= 'Z' {
			result = append(result, '-')
		}
		result = append(result, rune(strings.ToLower(string(c))[0]))
	}
	return string(result)
}
