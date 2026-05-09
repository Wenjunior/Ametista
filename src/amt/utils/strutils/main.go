package strutils

import (
	"github.com/dlclark/regexp2/v2"
	"github.com/dlclark/regexp2/v2/compat"
)

func FindAll(text string, expression string) []string {
	pattern := compat.MustCompile(expression, regexp2.RE2)

	matches := pattern.FindAllString(text, -1)

	return matches
}

func RemoveDuplicated(items []string) []string {
	keys := make(map[string] bool)

	result := []string{}

	for _, item := range items {
		value := keys[item]

		if !value {
			keys[item] = true

			result = append(result, item)
		}
	}

	return result
}

func Retain(items []string, expression string) []string {
	pattern := compat.MustCompile(expression, regexp2.RE2)

	var result []string

	for _, item := range items {
		if pattern.MatchString(item) {
			result = append(result, item)
		}
	}

	return result
}