package motivation

import (
	"math/rand"
	"os"
	"strings"
	"tg_ics_useful_bot/lib/e"
)

const quotesFile = "lib/motivation/quotes.txt"

var quotes = []string{}

func Quote() (string, error) {
	if len(quotes) != 0 {
		return quotes[rand.Intn(len(quotes))], nil
	}
	b, err := os.ReadFile(quotesFile)
	if err != nil {
		return "", e.Wrap("can't open quotes file", err)
	}
	strs := strings.Split(string(b), "\n")
	n := rand.Intn(len(strs))
	quotes = strs
	return strs[n], nil
}
