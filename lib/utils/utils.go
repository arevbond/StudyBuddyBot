package utils

import "strings"

func StringContains(x, str string) bool {
	for _, ch := range str {
		if string(ch) == x {
			return true
		}
	}
	return false
}

func Contains(strs []string, text string) bool {
	for _, str := range strs {
		if str == text {
			return true
		}
	}
	return false
}

func IsCommand(text string) bool {
	return strings.HasPrefix(text, "/")
}
