package quiz

import (
	"strings"
)

type Question struct {
	Question string
	Picture  string
	Answers  []string
}

func (q Question) IsCorrect(inAnswer string) bool {
	for _, answer := range q.Answers {
		if strings.TrimSpace(strings.ToLower(inAnswer)) == strings.TrimSpace(strings.ToLower(answer)) {
			return true
		}
	}
	return false
}
