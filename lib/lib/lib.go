package lib

func Contains(x, str string) bool {
	for _, ch := range str {
		if string(ch) == x {
			return true
		}
	}
	return false
}
