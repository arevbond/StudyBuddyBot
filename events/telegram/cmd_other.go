package telegram

import (
	"errors"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/lib/anecdots"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/flip"
	"tg_ics_useful_bot/lib/motivation"
	"tg_ics_useful_bot/storage"
)

// xkcdExec предоставляет Exec метод для выполнения /xkcd.
type xkcdExec string

// Exec: /xkcd - возвращает случайный xkcd комикс.
func (x xkcdExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	var comics xkcd.Comics
	comics, err := xkcd.RandomComics()
	if err != nil {
		return nil, e.Wrap("can't get comics from xkcd: ", err)
	}
	message := comics.Img
	mthd := sendPhotoMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// anekdotExec предоставляет Exec метод для выполнения /joke.
type anekdotExec string

func (a anekdotExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := anecdots.RandomAnecdot()
	if err != nil {
		return nil, e.Wrap("can't get anecdot: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

type addAnecdotExec string

func (a addAnecdotExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/add_anecdot)", errors.New("can't do this cmd"))
	}

	return nil, nil
}

// flipExec предоставляет Exec метод длы выполнения /flip.
type flipExec string

// Exec: /flip - возвращает случайную картинку из двух предоставленных ниже.
func (f flipExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := flip.KhinkalnyaOrVSU()
	mthd := sendPhotoMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// aufExec предоставляет Exec метод для выполнения /auf.
type aufExec string

// Exec: /auf - возвращает случайную мотивационную цитату.
func (a aufExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := motivation.Quote()
	if err != nil {
		return nil, e.Wrap("can't get quote: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}
