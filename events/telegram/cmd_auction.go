package telegram

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/lib/e"
	"tg_ics_useful_bot/storage"
	"time"
)

const (
	MaxDeposit     = 1000
	DefaultTimeout = 30
)

type AuctionPlayer struct {
	u       *storage.DBUser
	deposit int
}

// startAuctionExec предоставляет метод Exec для начала аукциона в чате.
type startAuctionExec string

// Exec: /start_auction [timeout] - запускает аукцион в чате, в котором указана данная команда.
// timeout - необязательный параметр времени окончания аукциона.
func (a startAuctionExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	if _, ok := p.auctions[chat.ID]; ok {
		return &Response{message: msgAuctionIsStarted, method: sendMessageMethod}, nil
	}

	p.auctions[chat.ID] = make([]*AuctionPlayer, 0)

	timeout := getAuctionTimeout(inMessage)
	go func() {
		time.Sleep(time.Duration(timeout) * time.Second)
		msg, err := a.finishAuction(chat.ID, p)
		if err != nil {
			log.Println("[ERROR] in goroutine /start_auction", err)
		}
		if msg != "" {
			_ = p.tg.SendMessage(chat.ID, msg, "", -1)
		}
	}()

	return &Response{message: fmt.Sprintf(msgStartAuction, MaxDeposit), method: sendMessageMethod, parseMode: telegram.Markdown}, nil
}

func getAuctionTimeout(inMessage string) int {
	strs := strings.Fields(inMessage)
	if len(strs) < 2 {
		return DefaultTimeout
	}
	timeout, err := strconv.Atoi(strs[1])
	if err != nil {
		return DefaultTimeout
	}
	if timeout < 10 || timeout > 60 {
		return DefaultTimeout
	}
	return timeout
}

// addDeposit предоставляет метод Exec для внесения депозита в ауцион.
type addDepositExec string

// Exec: /deposit {amount} - вносит депозит в текущий аукцион. Amount - обязательный параметр.
func (a addDepositExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message, err := a.addDeposit(inMessage, user, chat, p)
	if err != nil {
		return nil, e.Wrap("can't exec /deposit", err)
	}
	mthd := sendMessageMethod

	return &Response{message: message, method: mthd, replyMessageId: messageID}, nil
}

// addDeposit возвращает сообщание для телеграм чата, после команды /deposit {amount}.
func (a addDepositExec) addDeposit(inMessage string, user *telegram.User, chat *telegram.Chat, p *Processor) (string, error) {
	dbUser, err := p.storage.GetUser(context.Background(), user.ID, chat.ID)
	if err != nil {
		return "", err
	}

	if _, ok := p.auctions[chat.ID]; !ok {
		return msgAuctionNotStarted, nil
	}

	strs := strings.Fields(inMessage)
	if len(strs) < 2 {
		return msgErrorDepositCmd, nil
	}

	deposit, err := strconv.Atoi(strs[1])
	if err != nil {
		return msgErrorDepositCmd, nil
	}

	var needAdd bool
	player := getPlayer(dbUser, p.auctions)
	if player == nil {
		player = &AuctionPlayer{
			u:       dbUser,
			deposit: 0,
		}
		needAdd = true
	}

	if !a.canDeposit(deposit, dbUser, player) {
		return msgErrorDeposit, nil
	}

	if needAdd {
		p.auctions[chat.ID] = append(p.auctions[chat.ID], player)
	}
	err = p.changeDickSize(dbUser, -deposit)
	if err != nil {
		return "", err
	}
	player.deposit += deposit

	return fmt.Sprintf(msgSuccessDeposit, deposit), nil
}

// canDeposit проверяет может ли участник положить столько см пениса в аукцион.
func (a addDepositExec) canDeposit(deposit int, user *storage.DBUser, player *AuctionPlayer) bool {
	dickSize := user.DickSize
	playerDeposit := player.deposit
	return deposit >= 1 && deposit+playerDeposit <= MaxDeposit && dickSize-deposit >= 1
}

// getPlayer возвращает игрока аукциона.
func getPlayer(user *storage.DBUser, auctions map[int][]*AuctionPlayer) *AuctionPlayer {
	chatID := user.ChatID
	players := auctions[chatID]

	for _, p := range players {
		if user.ID == p.u.ID {
			return p
		}
	}

	return nil
}

// finishAuction случайным образом выбирает победителя из всех игроков аукциона.
func (a startAuctionExec) finishAuction(chatID int, p *Processor) (string, error) {
	if _, ok := p.auctions[chatID]; !ok {
		return msgAuctionNotStarted, nil
	}

	players := p.auctions[chatID]
	if len(players) == 0 {
		delete(p.auctions, chatID)
		return msgNotEnoughPlayers, nil
	}

	winner, reward := a.getAuctionWinnerAndReward(p.auctions[chatID])
	delete(p.auctions, chatID)

	for i := 3; i > 0; i-- {
		p.tg.SendMessage(chatID, fmt.Sprintf("До результата аукциона: %d!", i), "", -1)
		time.Sleep(1 * time.Second)
	}

	err := p.changeDickSize(winner, reward)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(msgWinner, winner.Username, reward), nil
}

// getAuctionWinnerAndReward случайным образом определяет победителя аукциона.
// Возвращает пользователя победителя и общий призовой фонд.
func (a startAuctionExec) getAuctionWinnerAndReward(players []*AuctionPlayer) (*storage.DBUser, int) {
	var reward int

	ids := make([]int, 0)

	for _, player := range players {
		reward += player.deposit
		for i := 1; i <= player.deposit; i++ {
			ids = append(ids, player.u.ID)
		}
	}

	winnerID := ids[rand.Intn(len(ids))]

	var winner *storage.DBUser
	for _, player := range players {
		if player.u.ID == winnerID {
			winner = player.u
		}
	}

	return winner, reward
}

// auctionExec предоставляет метод Exec для просмотра аукциона в чате.
type auctionExec string

// Exec: /auction - возвращает список всех участников текущего аукциона.
func (a auctionExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	var message string
	mthd := sendMessageMethod
	parseMode := telegram.Markdown

	if _, ok := p.auctions[chat.ID]; !ok {
		return &Response{message: msgAuctionNotStarted, method: mthd, parseMode: parseMode}, nil
	}

	message = a.getAuctionPlayers(chat.ID, p.auctions)

	return &Response{message: message, method: mthd, parseMode: parseMode}, nil
}

// getAuctionPlayers возвращает список текущих участников аукциона.
func (a auctionExec) getAuctionPlayers(chatID int, auctions map[int][]*AuctionPlayer) string {
	players := auctions[chatID]

	if len(players) == 0 {
		return msgZeroPlayers
	}

	message := "Текущие игроки аукциона:\n\n"
	reward := 0

	for _, p := range players {
		if p.deposit > 0 {
			reward += p.deposit
			message += fmt.Sprintf("%s:\n*8", p.u.FirstName+" "+p.u.LastName)
			for i := 0; i < p.deposit/5; i++ {
				message += "="
			}
			message += "=Ð*\n"
		}
	}
	message += fmt.Sprintf("\nОбщий фонд *%d см*!", reward)
	return message
}
