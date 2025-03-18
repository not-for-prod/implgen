package strtools

import "strings"

// SnakeCase converts camelCase or PascalCase to snake_case
func SnakeCase(str string) string {
	var result []rune
	prevLower := false

	for i, c := range str {
		if i > 0 && ('A' <= c && c <= 'Z') {
			// Only insert underscore if the previous character is lowercase or next is lowercase (handling InvoiceP2P case)
			if prevLower || (i+1 < len(str) && 'a' <= rune(str[i+1]) && rune(str[i+1]) <= 'z') {
				result = append(result, '_')
			}
		}
		result = append(result, rune(strings.ToLower(string(c))[0]))
		prevLower = ('a' <= c && c <= 'z')
	}

	return string(result)
}

func KebabCase(str string) string {
	var result []rune
	prevLower := false

	for i, c := range str {
		if i > 0 && ('A' <= c && c <= 'Z') {
			// Only insert underscore if the previous character is lowercase or next is lowercase (handling InvoiceP2P case)
			if prevLower || (i+1 < len(str) && 'a' <= rune(str[i+1]) && rune(str[i+1]) <= 'z') {
				result = append(result, '-')
			}
		}
		result = append(result, rune(strings.ToLower(string(c))[0]))
		prevLower = ('a' <= c && c <= 'z')
	}

	return string(result)
}
