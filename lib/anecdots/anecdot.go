package anecdots

import (
	"math/rand"
	"os"
	"tg_ics_useful_bot/lib/e"

	"gopkg.in/yaml.v3"
)

var anecdotsFileName = "lib/anecdots/anecdots.yaml"

type anecdot struct {
	Text string `yaml:"text"`
}

var anecdots = []anecdot{}

func RandomAnecdot() (string, error) {
	if len(anecdots) != 0 {
		return anecdots[rand.Intn(len(anecdots))].Text, nil
	}
	b, err := os.ReadFile(anecdotsFileName)
	if err != nil {
		return "", e.Wrap("can't find anecdots file", err)
	}
	err = yaml.Unmarshal(b, &anecdots)
	if err != nil {
		return "", e.Wrap("can't unmarshall anecdots file", err)
	}
	return anecdots[rand.Intn(len(anecdots))].Text, nil
}
