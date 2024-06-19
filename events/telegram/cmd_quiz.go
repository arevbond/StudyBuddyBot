package telegram

// TODO: add index in /quit namequiz.yaml {index_question}
// TODO: add /pause && /continue
import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sort"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/lib/quiz"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	defaultAward         = 1
	timeBetweenQuestions = 5
)

type startQuizExec string

func (s startQuizExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/star_quiz)", errors.New("can't do this cmd"))
	}

	strs := strings.Split(inMessage, " ")
	if len(strs) != 2 {
		return &Response{message: msgWriteQuizName, method: sendMessageMethod}, nil
	}
	filename := strs[1]

	quizGame, err := quiz.New(filename)

	if err != nil {
		_ = p.tg.SendMessage(chat.ID, fmt.Sprintf(msgErrorQuiz, filename), telegram.WithoutParseMode, messageID)
		return nil, e.Wrap("can't start quit", err)
	}

	go s.startQuiz(quizGame, chat.ID, p)

	return &Response{message: fmt.Sprintf(msgStartQuiz, quizGame.Theme, quizGame.GetLevel(), len(quizGame.Questions)), method: sendMessageMethod,
		parseMode: telegram.Markdown}, nil
}

func (s startQuizExec) startQuiz(quizGame quiz.Quiz, chatID int, p *Processor) {
	time.Sleep(5 * time.Second)

	for _, question := range quizGame.Questions {
		select {
		case <-p.quiz.quit:
			p.quiz.currentPlayers = make(map[int]int)
			p.quiz.currentQuestion = &quiz.Question{}
			return
		default:
			p.quiz.currentQuestion = question
			if question.Picture != "" {
				_ = p.tg.SendPhoto(chatID, question.Picture)
			}
			_ = p.tg.SendPoll(telegram.NewSendPoll(chatID, question))
			time.Sleep(time.Duration(question.OpenPeriod+timeBetweenQuestions) * time.Second)
		}
	}

	awardMessage := s.awarding(chatID, quizGame.Level, p.quiz.currentPlayers, p)
	_ = p.tg.SendMessage(chatID, msgFinishQuiz+"\n"+awardMessage, "", -1)
	p.quiz.currentPlayers = make(map[int]int)
	p.quiz.currentQuestion = &quiz.Question{}
}

func (s startQuizExec) awarding(chatID int, level quiz.Level, players map[int]int, p *Processor) string {
	award := defaultAward

	sortedPlayers := getSortedQuizPlayers(players)
	result := "Результаты:\n"

	for _, player := range sortedPlayers {
		dbUser, err := p.storage.GetUser(context.Background(), player, chatID)
		if err != nil {
			p.logger.Error("can't get db user", slog.String("func", "startQuizExec.awarding"))
			continue
		}
		dbUser.DickSize += p.quiz.currentPlayers[player] * award
		err = p.storage.UpdateUser(context.Background(), dbUser)
		if err != nil {
			p.logger.Error("can't update points in db use", slog.String("func", "startQuizExec.awarding"))
			continue
		}
		result += fmt.Sprintf(" • %d ✔  %s          ➕ %d см\n", p.quiz.currentPlayers[player], dbUser.FirstName+" "+dbUser.LastName,
			p.quiz.currentPlayers[player]*award)
	}
	return result
}

func getSortedQuizPlayers(currentPlayers map[int]int) []int {
	players := []int{}
	for player, _ := range currentPlayers {
		players = append(players, player)
	}
	sort.Slice(players, func(i, j int) bool {
		return currentPlayers[players[i]] > currentPlayers[players[j]]
	})
	return players
}

type stopQuizExec string

func (s stopQuizExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if !p.isAdmin(user.ID) {
		return nil, e.Wrap("no admin can't do this cmd (/star_quiz)", errors.New("can't do this cmd"))
	}

	if len(p.quiz.currentQuestion.Options) > 0 {
		p.quiz.quit <- true
		return &Response{message: msgStoppedQuiz, method: sendMessageMethod, replyMessageId: messageID}, nil
	}
	return &Response{message: msgQuizNotAvailable, method: sendMessageMethod, replyMessageId: messageID}, nil
}
