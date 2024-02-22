package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/quiz"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	award                = 250
	timeBetweenQuestions = 5
)

var currentQuestion = &quiz.Question{}
var currentPlayers = make(map[int]int)

type startQuizExec string

func (s startQuizExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/star_quiz)", errors.New("can't do this cmd"))
	}

	strs := strings.Split(inMessage, " ")
	if len(strs) != 2 {
		return &Response{message: "Введите название quiz", method: sendMessageMethod}, nil
	}
	filename := strs[1]

	quizGame, err := quiz.New(filename)

	if err != nil {
		return nil, e.Wrap("can't start quiz", err)
	}

	go p.startQuiz(quizGame.Questions, chat.ID)

	return &Response{message: fmt.Sprintf(msgStartQuiz, quizGame.Theme, quizGame.GetLevel(), len(quizGame.Questions)), method: sendMessageMethod,
		parseMode: telegram.Markdown}, nil
}

func (p *Processor) startQuiz(questions []*quiz.Question, chatID int) {
	time.Sleep(10 * time.Second)

	for _, question := range questions {
		currentQuestion = question
		if question.Picture != "" {
			_ = p.tg.SendPhoto(chatID, question.Picture)
		}
		_ = p.tg.SendPoll(telegram.NewSendPoll(chatID, question))
		time.Sleep(time.Duration(question.OpenPeriod+timeBetweenQuestions) * time.Second)
	}

	awardMessage := p.awarding(chatID)
	_ = p.tg.SendMessage(chatID, msgFinishQuiz+"\n"+awardMessage, "", -1)
	currentPlayers = make(map[int]int)
	currentQuestion = &quiz.Question{}
}

func (p *Processor) awarding(chatID int) string {
	players := []int{}
	for player, _ := range currentPlayers {
		players = append(players, player)
	}
	sort.Slice(players, func(i, j int) bool {
		return currentPlayers[players[i]] > currentPlayers[players[j]]
	})

	result := "Результаты:\n"

	for _, player := range players {
		dbUser, err := p.storage.GetUser(context.Background(), player, chatID)
		if err != nil {
			log.Println("can't get db user", err)
			continue
		}
		dbUser.DickSize += currentPlayers[player] * award
		err = p.storage.UpdateUser(context.Background(), dbUser)
		if err != nil {
			log.Println("can't update points in db user", err)
			continue
		}
		result += fmt.Sprintf("%s: %d п. о. ➕ %d см\n", dbUser.FirstName+" "+dbUser.LastName, currentPlayers[player],
			currentPlayers[player]*award)
	}
	return result
}
