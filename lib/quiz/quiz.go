package quiz

import (
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"tg_ics_useful_bot/lib/e"
)

const pathToFile = "lib/quiz/quizzes/"

type level int

const (
	Easy   level = 1
	Medium       = 2
	Hard         = 3
)

const (
	defaultOpenPeriod = 20
)

type Quiz struct {
	Theme     string      `json:"theme" yaml:"theme"`
	Level     level       `json:"level" yaml:"level,omitempty"`
	Questions []*Question `json:"questions" yaml:"questions"`
}

func (q Quiz) GetLevel() string {
	switch q.Level {
	case Easy:
		return "Easy"
	case Medium:
		return "Medium"
	case Hard:
		return "Hard"
	}
	return "Unknown"
}

func New(filename string) (Quiz, error) {
	data, err := os.ReadFile(pathToFile + filename)
	if err != nil {
		return Quiz{}, e.Wrap("can't read file", err)
	}
	var quiz Quiz

	switch strings.Split(filename, ".")[1] {
	case "json":
		err := json.Unmarshal(data, &quiz)
		if err != nil {
			return Quiz{}, e.Wrap("can't unmarshall json", err)
		}
	case "yaml":
		err := yaml.Unmarshal(data, &quiz)
		if err != nil {
			return Quiz{}, e.Wrap("can't unmarshall json", err)
		}
	}

	n := len(quiz.Questions)

	for i, q := range quiz.Questions {
		if q.OpenPeriod < 5 {
			q.OpenPeriod = defaultOpenPeriod
		}
		q.Question += fmt.Sprintf(" [%d/%d]", i+1, n)
	}

	return quiz, nil
}

type Question struct {
	Question              string   `json:"question" yaml:"question"`
	Picture               string   `json:"picture,omitempty" yaml:"picture,omitempty"`
	Options               []string `json:"options" yaml:"options"`
	CorrectOptionID       int      `json:"correct_option_id" yaml:"correct_option_id"`
	AllowsMultipleAnswers bool     `json:"allows_multiple_answers" yaml:"allows_multiple_answers"`
	Explanation           string   `json:"explanation" yaml:"explanation"`
	OpenPeriod            int      `default:"15" json:"open_period" yaml:"open_period"`
}

func (q Question) IsCorrect(message string) bool {
	for _, answer := range q.Options {
		if strings.TrimSpace(strings.ToLower(message)) == strings.TrimSpace(strings.ToLower(answer)) {
			return true
		}
	}
	return false
}
