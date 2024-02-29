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
	defaultBonus = 50
)

func (l Level) Bonus() int {
	switch l {
	case Easy, Medium, Hard, VeryHard:
		return int(l) * defaultBonus
	}
	return 0
}

const (
	Easy     Level = 1
	Medium         = 2
	Hard           = 3
	VeryHard       = 4
)

const (
	maxLenQuestion    = 300
	maxLenExplanation = 200
	maxLenOption      = 100
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
	formatFields(quiz.Questions)
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
		q.Question += fmt.Sprintf(" [%d/%d]", i+1, n)
	}
}

func formatFields(questions []*Question) {
	for i, q := range questions {
		if len([]rune(q.Question)) > maxLenQuestion {
			questions[i].Question = format(q.Question, maxLenQuestion)
		}
		if len([]rune(q.Explanation)) > maxLenExplanation {
			questions[i].Explanation = format(q.Explanation, maxLenExplanation)
		}
		for j, option := range q.Options {
			if len([]rune(option)) > maxLenOption {
				questions[i].Options[j] = format(option, maxLenOption)
			}
		}
	}
}

func format(str string, maxLen int) string {
	chars := strings.Split(str, "")
	return strings.Join(chars[:maxLen], "")
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
