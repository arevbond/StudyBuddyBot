package utils

import (
	"regexp"
	"strings"
	"unicode"
)

const (
	UnsupportedCommand = iota
	IsYesCommand
	IsNoCommand
)

func CheckYesOrNo(text string) int {
	text = strings.TrimSpace(text)
	var isWord = regexp.MustCompile(`(\w|[а-яА-Я])`)

	isYes := func(word string) bool {
		var validString = regexp.MustCompile(`(^(Д|д)+(а|a|A|А)+\p{P}+$)|(^(Д|д)+(а|a|A|А)+$)`)
		return validString.MatchString(word)
	}
	isNo := func(word string) bool {
		var validString = regexp.MustCompile(`(^((н|Н)+(e|E|е|Е)+(т|Т)+\p{P}+$))|(^(н|Н)+(e|E|е|Е)+(т|Т)+$)`)
		return validString.MatchString(word)
	}

	textSplited := strings.Split(text, " ")

	words := make([]string, 0)

	for _, w := range textSplited {
		if isWord.MatchString(w) {
			words = append(words, w)
		}
	}

	if len(words) == 0 {
		return UnsupportedCommand
	}

	if isYes(words[len(words)-1]) {
		return IsYesCommand
	}

	if isNo(words[len(words)-1]) {
		return IsNoCommand
	}
	return UnsupportedCommand
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
