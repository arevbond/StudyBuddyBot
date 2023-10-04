package utils

import (
	"strings"
	"unicode"
)

func IsYesCommand(text string) bool {
	var flag bool
	text = strings.ToLower(text)
	text = removeSpecSymbols(text)
	if strings.HasPrefix(text, "д") {
		for _, ch := range []rune(text) {
			if ch != 'д' && ch != 'а' {
				//log.Printf("char: %q != 'д' и 'а'\n")
				return false
			}
			if ch != 'д' && flag == false {
				flag = true
			} else if ch == 'д' && flag != false {
				return false
			} else if ch == 'а' && flag || ch == 'д' {
				continue
			} else {
				//log.Printf("%q - почему то вернуло false", ch)
				return false
			}
		}
	} else {
		return false
	}
	return true
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
