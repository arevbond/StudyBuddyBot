package telegram

import (
	"log/slog"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
)

const (
	defaultPrompt = "привет!"

	scriptPath = "lib/groq/script.py"
	venvPath   = "lib/groq/venv/bin/python3"
)

func answerFromLlama3(p *Processor, chatID int, inMessage string, messageID int) error {

	prompt := parsePrompt(inMessage)

	cmd := exec.Command(venvPath, scriptPath, prompt)

	cmd.Env = append(os.Environ(), "PYTHONPATH=lib/groq/")

	output, err := cmd.Output()
	if err != nil {
		p.logger.Error("can't get output from python script", slog.Any("error", err))
		return e.Wrap("can't get output from python script", err)
	}
	return p.tg.SendMessage(chatID, string(output), telegram.WithoutParseMode, messageID)
}

func parsePrompt(inMessage string) string {
	strs := strings.Fields(inMessage)
	if len(strs) > 1 {
		return strings.Join(strs[1:], " ")
	}
	return defaultPrompt
}

func isAppeal(text string) bool {
	words := strings.Fields(text)
	var validAppeal = regexp.MustCompile(`(?i)^(аркадий|аркаш|бот)`)
	return len(words) > 0 && validAppeal.MatchString(words[0])
}
