package telegram

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/quiz"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	timeToAnswer = 20
	totalAward   = 5000
)

var isAnswered bool
var chatToCurrentQuestion = make(map[int]quiz.Question)
var chatToPlayers = make(map[int]map[int]int)

type startQuizExec string

func (s startQuizExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/star_quiz)", errors.New("can't do this cmd"))
	}

	if len(quiz.Quizzes) < 1 {
		return &Response{message: msgZeroQuizzes, replyMessageId: messageID, method: sendMessageMethod}, nil
	}

	var number int = 0

	strs := strings.Split(inMessage, " ")
	if len(strs) >= 2 {
		num, err := strconv.Atoi(strs[1])
		if err == nil && len(quiz.Quizzes) <= num {
			number = num - 1
		}
	}

	quiz := quiz.Quizzes[number]

	go p.startQuiz(quiz.Questions, chat.ID)

	return &Response{message: fmt.Sprintf(msgStartQuiz, quiz.Theme, len(quiz.Questions)), method: sendMessageMethod,
		parseMode: telegram.Markdown}, nil
}

func (p *Processor) startQuiz(questions []quiz.Question, chatID int) {
	chatToPlayers[chatID] = make(map[int]int)

	time.Sleep(5 * time.Second)

	for i, question := range questions {
		isAnswered = false
		message := fmt.Sprintf("Вопрос №%d\n", i+1)
		p.tg.SendMessage(chatID, message+question.Question, "", -1)
		if question.Picture != "" {
			p.tg.SendPhoto(chatID, question.Picture)
		}
		chatToCurrentQuestion[chatID] = question
		time.Sleep(timeToAnswer * time.Second)
	}

	awardMessage := p.awarding(chatID)
	p.tg.SendMessage(chatID, msgFinishQuiz+"\n"+awardMessage, "", -1)
	delete(chatToPlayers, chatID)
	delete(chatToCurrentQuestion, chatID)
}

func (p *Processor) checkAnswer(chatID int, tgID int, answer string, messageID int) {
	question := chatToCurrentQuestion[chatID]
	if question.IsCorrect(answer) {
		chatToPlayers[chatID][tgID]++
		if !isAnswered {
			isAnswered = true
			p.tg.SendMessage(chatID, "👍", "", messageID)
		}
	}
}

func (p *Processor) awarding(chatID int) string {
	playersToScore := chatToPlayers[chatID]
	if playersToScore == nil {
		return ""
	}

	players := []int{}
	for player, _ := range playersToScore {
		players = append(players, player)
	}
	sort.Slice(players, func(i, j int) bool {
		return playersToScore[players[i]] > playersToScore[players[j]]
	})

	result := "Результат:\n"

	award := totalAward

	for _, player := range players {
		dbUser, err := p.storage.GetUser(context.Background(), player, chatID)
		if err != nil {
			log.Println("can't get db user", err)
			continue
		}
		dbUser.Points += award
		err = p.storage.UpdateUser(context.Background(), dbUser)
		if err != nil {
			log.Println("can't update points in db user", err)
			continue
		}
		result += fmt.Sprintf("%s: %d + %d очков\n", dbUser.FirstName+" "+dbUser.LastName, playersToScore[player], award)
		award = award / 2
	}
	return result
}
