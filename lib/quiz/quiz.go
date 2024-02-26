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

type Level int

const (
	Easy     Level = 1
	Medium         = 2
	Hard           = 3
	VeryHard       = 4
)

const (
	defaultOpenPeriod = 20
)

type Quiz struct {
	Theme     string      `json:"theme" yaml:"theme"`
	Level     Level       `json:"level" yaml:"level,omitempty"`
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
	case VeryHard:
		return "VeryHard"
	}
	return "Unknown"
}

func New(filename string) (Quiz, error) {
	quiz, err := readQuizFromFile(filename)
	if err != nil {
		return Quiz{}, e.Wrap("can't read quiz", err)
	}
	addIndexes(quiz.Questions)
	return quiz, nil
}

func readQuizFromFile(filename string) (Quiz, error) {
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
			return Quiz{}, e.Wrap("can't unmarshall yaml", err)
		}
	}
	return quiz, nil
}

func addIndexes(questions []*Question) {
	n := len(questions)

	for i, q := range questions {
		if q.OpenPeriod < 5 {
			q.OpenPeriod = defaultOpenPeriod
		}
		q.Question += fmt.Sprintf(" [%d/%d]", i+1, n)
	}
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
