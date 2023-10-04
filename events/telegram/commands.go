package telegram

import (
	"fmt"
	"log"
	"strings"
	"tg_ics_useful_bot/clients/jokesrv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/lessons"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/lib"
)

func (p *Processor) do(method method, chatID int, message string) error {
	switch method {
	case sendMessageMethod:
		return p.tg.SendMessage(chatID, message)
	case sendPhotoMethod:
		return p.tg.SendPhoto(chatID, message)
	}
	return e.Wrap(fmt.Sprintf("unsupported method: %q", message), nil)
}

func (p *Processor) doCmd(text string, chat *telegram.Chat, user *telegram.User, messageID int) error {
	text = strings.TrimSpace(text)
	splitedText := strings.Split(text, " ")
	if lib.IsYes(splitedText[len(splitedText)-1]) {
		return p.tg.SendMessage(chat.ID, "Пизда")
	}
	if strings.HasPrefix(text, "/") {
		log.Printf("[INFO] got new command '%s' from '%s' in %s", text, user.Username, chat.Title)
	}
	if chat.Type == "group" || chat.Type == "supergroup" {
		switch {
		case strings.HasPrefix(text, AllCmd):
			message := p.AllUsernames(chat.ID)
			return p.do(sendMessageMethod, chat.ID, message)

		case strings.HasPrefix(text, GayTopCmd):
			message, err := p.gameGayTop(chat.ID)
			if err != nil {
				return e.Wrap("can't do GayTopCmd: ", err)
			}
			return p.do(sendMessageMethod, chat.ID, message)
		case strings.HasPrefix(text, GayStartCmd):
			message, err := p.gameGay(chat.ID)
			if err != nil {
				return e.Wrap("can't get message from gameGay: ", err)
			}
			return p.do(sendMessageMethod, chat.ID, message)

		case strings.HasPrefix(text, DickTopCmd):
			message, err := p.topDick(chat.ID)
			if err != nil {
				return e.Wrap(fmt.Sprintf("can't get top dics from chat %d: ", chat.ID), err)
			}
			return p.do(sendMessageMethod, chat.ID, message)
		case strings.HasPrefix(text, DicStartCmd):
			message, err := p.gameDick(chat, user, messageID)
			if err != nil {
				return e.Wrap("can't get message from gameDick: ", err)
			}
			return p.do(sendMessageMethod, chat.ID, message)
		case strings.HasPrefix(text, DickDuelCmd):
			message, err := p.gameDuelDick(chat, messageID, user, user.Username)
			if err != nil {
				return e.Wrap("can't do gameDuelDick: ", err)
			}
			if lib.Contains("@", text) {
				textSplited := strings.Split(text, "@")
				target := textSplited[len(textSplited)-1]
				log.Printf("[INFO] @%s вызывает на дуель @%s", user.Username, target)
				message, err = p.gameDuelDick(chat, messageID, user, target)
				if err != nil {
					return e.Wrap("can't do gameDuelDick: ", err)
				}
			}
			return p.do(sendMessageMethod, chat.ID, message)

		case strings.HasPrefix(text, TodayLessonsCmd):
			message := lessons.LessonsToday()
			return p.do(sendMessageMethod, chat.ID, message)
		case strings.HasPrefix(text, TomorrowLessonsCmd):
			message := lessons.TomorrowLessons()
			return p.do(sendMessageMethod, chat.ID, message)
		case strings.HasPrefix(text, LessonsCmd):
			message := lessons.AllLessons()
			return p.do(sendMessageMethod, chat.ID, message)

		case strings.HasPrefix(text, XkcdCmd):
			comics, err := xkcd.RandomComics()
			if err != nil {
				return e.Wrap("can't get comics from xkcd: ", err)
			}
			return p.do(sendPhotoMethod, chat.ID, comics.Img)
		case strings.HasPrefix(text, AnecdotCmd):
			anecdot, err := jokesrv.Anecdot()
			if err != nil {
				return e.Wrap("can't get anecdot: ", err)
			}
			return p.do(sendMessageMethod, chat.ID, anecdot)
		case strings.HasPrefix(text, FlipCmd):
			return p.do(sendPhotoMethod, chat.ID, RandomPhotoHinkOrRoom())
		default:
			return nil
		}
	}
	return nil
}

func (p *Processor) AllUsernames(chatID int) string {
	admins, err := p.tg.ChatAdministrators(chatID)
	if err != nil {
		log.Printf("can't get admins in chat #%d: ", chatID, err)
	}
	result := ""
	for _, a := range admins {
		result += "@" + a.Username + " "
	}
	return result[:len(result)-1]
}
