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

func Equal(text string, strs []string) bool {
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

func Abs(a int) int {
	if a > 0 {
		return a
	}
	return -1 * a
}
