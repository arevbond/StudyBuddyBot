package utils

import (
	"regexp"
	"unicode"
)

func IsYesCommand(text string) bool {
	var validString = regexp.MustCompile(`\B(Д|д)+(а|a|A|А)+\B`)
	return validString.MatchString(text)
}

func IsNoCommand(text string) bool {
	var validString = regexp.MustCompile(`\B(н|Н)+(e|E|е|Е)+(т|Т)+\B`)
	return validString.MatchString(text)
}

func removeSpecSymbols(text string) string {
	result := ""
	for _, ch := range text {
		if unicode.IsLetter(rune(ch)) {
			result += string(ch)
		}
	}
	return result
}
