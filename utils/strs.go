package utils

import "strings"

func ToTitle(underscoreName string) string {
	words := strings.Split(strings.ToLower(underscoreName), "_")
	for i := 0; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}
func ToCamelCase(underscoreName string) string {
	words := strings.Split(strings.ToLower(underscoreName), "_")
	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}
	return strings.Join(words, "")
}
