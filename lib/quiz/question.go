package quiz

import (
	"strings"
	"time"
)

type Question struct {
	Question     string
	Picture      string
	Answers      []string
	TimeToAnswer time.Duration
}

func (q Question) IsCorrect(inAnswer string) bool {
	for _, answer := range q.Answers {
		if strings.TrimSpace(strings.ToLower(inAnswer)) == strings.TrimSpace(strings.ToLower(answer)) {
			return true
		}
	}
	return false
}
