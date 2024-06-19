package telegram

import (
	"log/slog"
	"os"
	"os/exec"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
)

type llamaExec string

func (l llamaExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	prompt := l.parsePrompt(inMessage)

	scriptPath := "lib/groq/script.py"

	venvPath := "lib/groq/venv/bin/python3"

	cmd := exec.Command(venvPath, scriptPath, prompt)

	cmd.Env = append(os.Environ(), "PYTHONPATH=lib/groq/")

	output, err := cmd.Output()
	if err != nil {
		p.logger.Error("can't get output from python script", slog.Any("error", err))
		return nil, e.Wrap("can't get output from python script", err)
	}
	return &Response{message: string(output), method: sendMessageMethod, replyMessageId: messageID}, nil
}

func (l llamaExec) parsePrompt(inMessage string) string {
	strs := strings.Fields(inMessage)
	if len(strs) > 1 {
		return strings.Join(strs[1:], " ")
	}
	return "Привет, Аркадий, можешь рассказать о себе?"
}
