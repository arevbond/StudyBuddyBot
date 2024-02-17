package quiz

import (
	"strings"
)

type Question struct {
	Question     string   `json:"question" yaml:"question"`
	Picture      string   `json:"picture,omitempty" yaml:"picture,omitempty"`
	Answers      []string `json:"answers" yaml:"answers"`
	TimeToAnswer int      `json:"time_to_answer,omitempty" yaml:"time_to_answer,omitempty"`
}

func (q Question) IsCorrect(message string) bool {
	for _, answer := range q.Answers {
		if strings.TrimSpace(strings.ToLower(message)) == strings.TrimSpace(strings.ToLower(answer)) {
			return true
		}
	}
	return false
}
