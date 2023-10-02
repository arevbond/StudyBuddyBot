package lib

import (
	"strings"
)

func IsYes(text string) bool {
	var flag bool
	text = strings.ToLower(text)
	if strings.HasPrefix(text, "д") {
		for _, ch := range []rune(text) {
			//log.Printf("%q\n", ch)
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
