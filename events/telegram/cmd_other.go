package telegram

import (
	"math/rand"
	"tg_ics_useful_bot/clients/jokesrv"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/clients/xkcd"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

// xkcdExec предоставляет Exec метод для выполнения /xkcd.
type xkcdExec string

// Exec: /xkcd - возвращает случайный xkcd комикс.
func (a xkcdExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	var comics xkcd.Comics
	comics, err := xkcd.RandomComics()
	if err != nil {
		return nil, e.Wrap("can't get comics from xkcd: ", err)
	}
	message := comics.Img
	mthd := sendPhotoMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// anekdotExec предоставляет Exec метод для выполнения /joke.
type anekdotExec string

// Exec: /joke - возвращает случайный анекдот от @bobuk.
func (a anekdotExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := jokesrv.Anecdot()
	if err != nil {
		return nil, e.Wrap("can't get anecdot: ", err)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// flipExec предоставляет Exec метод длы выполнения /flip.
type flipExec string

// Exec: /flip - возвращает случайную картинку из двух предоставленных ниже.
func (a flipExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := khinkalnyaOrVSU()
	mthd := sendPhotoMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

const (
	urlImgHink = "https://avatars.mds.yandex.net/get-altay/3518606/2a00000179e2472a99931c431d308fd69e09/XXL"
	urlImgRoom = "https://www.vsu.ru/gallery/photos/study/dept_phys.jpg"
)

// khinkalnyaOrVSU возвращает URL картини хинкальни или VSU аудитории.
func khinkalnyaOrVSU() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(2)
	if n == 1 {
		return urlImgHink
	}
	return urlImgRoom
}
