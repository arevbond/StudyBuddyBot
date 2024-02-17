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
	defaultTimeToAnswer = 75
	award               = 200
)

var isAnswered = make(map[int]bool) // TODO: refactor this var
var chatToCurrentQuestion = make(map[int]quiz.Question)
var chatToPlayers = make(map[int]map[int]int)

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

func (p *Processor) startQuiz(questions []quiz.Question, chatID int) {
	chatToPlayers[chatID] = make(map[int]int)

	time.Sleep(7 * time.Second)

	for i, question := range questions {
		isAnswered[chatID] = false
		chatToCurrentQuestion[chatID] = question
		message := fmt.Sprintf("Вопрос №%d\n", i+1)
		_ = p.tg.SendMessage(chatID, message+question.Question, "", -1)
		if question.Picture != "" {
			_ = p.tg.SendPhoto(chatID, question.Picture)
		}

		timeToAnswer := question.TimeToAnswer
		if timeToAnswer <= 0 {
			timeToAnswer = defaultTimeToAnswer
		}

		n := timeToAnswer
		for i := 0; i < n; i++ {
			if isAnswered[chatID] {
				break
			}
			time.Sleep(1 * time.Second)

		}

		if !isAnswered[chatID] {
			message = "Правильный ответ:\n"
			if len(question.Answers) > 0 {
				_ = p.tg.SendMessage(chatID, message+question.Answers[0], "", -1)
			}
		}
		isAnswered[chatID] = true
		time.Sleep(20 * time.Second)
	}

	awardMessage := p.awarding(chatID)
	_ = p.tg.SendMessage(chatID, msgFinishQuiz+"\n"+awardMessage, "", -1)
	delete(chatToPlayers, chatID)
	delete(chatToCurrentQuestion, chatID)
}

func (p *Processor) checkAnswer(chatID int, tgID int, answer string, messageID int) {
	question := chatToCurrentQuestion[chatID]
	if question.IsCorrect(answer) {
		if !isAnswered[chatID] {
			chatToPlayers[chatID][tgID]++
			isAnswered[chatID] = true
			_ = p.tg.SendMessage(chatID, "👍", "", messageID)
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

	result := "Результаты:\n"

	for _, player := range players {
		dbUser, err := p.storage.GetUser(context.Background(), player, chatID)
		if err != nil {
			log.Println("can't get db user", err)
			continue
		}
		dbUser.DickSize += playersToScore[player] * award
		err = p.storage.UpdateUser(context.Background(), dbUser)
		if err != nil {
			log.Println("can't update points in db user", err)
			continue
		}
		result += fmt.Sprintf("%s: %d п. о. ➕ %d см\n", dbUser.FirstName+" "+dbUser.LastName, playersToScore[player],
			playersToScore[player]*award)
	}
	return result
}
