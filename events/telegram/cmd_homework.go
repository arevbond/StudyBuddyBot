package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/storage"
)

type Homework struct {
	subject string
	Task    string
}

func newHomework(subject, task string) *Homework {
	return &Homework{subject: subject, Task: task}
}

const (
	maxRows = 5
)

type UserWithChat struct {
	ChatID int
	UserID int
}

var stateHomework = make(map[UserWithChat]*Homework)

type addHomeworkExec string

func (a addHomeworkExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {
	
	message := p.addHomeworkCmd(inMessage, UserWithChat{ChatID: chat.ID, UserID: user.ID})
	mthd := sendMessageWithButtonsMethod
	replyMessageId := messageID
	return &Response{message: message, method: mthd, replyMessageId: replyMessageId}, nil
}

func (p *Processor) addHomeworkCmd(text string, userWithChat UserWithChat) string {
	if strings.HasPrefix(text, "/") {
		stateHomework[userWithChat] = newHomework("", "")
		return msgAddSubject
	} else if hm, ok := stateHomework[userWithChat]; ok && hm.subject == "" {
		hm.subject = text
		return msgAddTask
	} else if hm, ok = stateHomework[userWithChat]; ok && hm.Task == "" {
		hm.Task = text
		message := fmt.Sprintf("ДЗ: %s - %s успешно добавлено", hm.subject, hm.Task)
		err := p.storage.AddHomework(context.Background(), userWithChat.ChatID, hm.subject, hm.Task)
		if err != nil {
			message = msgErrorAddHomework
			log.Printf("can't add homework: %v", err)
		}
		delete(stateHomework, userWithChat)
		return message
	}
	return msgSomethingWrong
}

// getHomeworkExec предоставляет метод Exec для выполнения /get.
type getHomeworkExec string

// Exec: /get [number] [subject] - возвращает последние записи домашнего задания
func (a getHomeworkExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	message := p.getHomework(inMessage, chat.ID)
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// getHomework формирует строку домашнего задания.
func (p *Processor) getHomework(text string, chatID int) string {
	val := ""
	for _, s := range strings.Split(text, " ")[1:] {
		if s != "" {
			_, err := strconv.Atoi(s)
			if err == nil {
				val = s
				break
			}
			val += s + " "
		}
	}

	message := ""

	homeworks := []*storage.DBHomework{}
	if num, err := strconv.Atoi(val); err == nil {
		homeworks, err = p.storage.GetHomeworkByChatID(context.Background(), chatID, num)
		if err != nil {
			log.Print(err)
			return ""
		}
		message += fmt.Sprintf("Последние %d домашних задания:\n", num)
	} else if val != "" {
		val = val[:len(val)-1]
		homeworks, err = p.storage.GetHomeworkBySubject(context.Background(), chatID, val)
		if err != nil {
			log.Print(err)
			return ""
		}
		message += fmt.Sprintf("Всё домашнее задание по предмету %s:\n", val)
	} else {
		homeworks, err = p.storage.GetHomeworkByChatID(context.Background(), chatID, maxRows)
		if err != nil {
			log.Print(err)
			return ""
		}
		message += fmt.Sprintf("Последние %d добавленных домашних задания:\n", maxRows)
	}

	for _, hm := range homeworks {
		message += fmt.Sprintf(" • \"%s\" - \"%s\". [id = %d]\n", hm.Subject, hm.Task, hm.ID)
	}

	return message
}

// deleteHomeworkExec предоставляет метод Exec для выполнения /delete.
type deleteHomeworkExec string

// Exec: /delete [id] - удаляет запись о домашнем задании
func (a deleteHomeworkExec) Exec(p *Processor, inMessage string, user *telegram.User, chat *telegram.Chat,
	userStats *storage.DBUserStat, messageID int) (*Response, error) {

	val := ""
	for _, str := range strings.Split(inMessage, " ")[1:] {
		if str != "" {
			val = str
			break
		}
	}
	num, err := strconv.Atoi(val)
	message := p.deleteHomework(num)
	if err != nil {
		message = fmt.Sprintf(msgIncorrectValue, val)
	}
	mthd := sendMessageMethod
	return &Response{message: message, method: mthd, replyMessageId: -1}, nil
}

// deleteHomework удаляет запись домашнего задания.
func (p *Processor) deleteHomework(rowID int) string {
	err := p.storage.DeleteHomework(context.Background(), rowID)
	message := fmt.Sprintf(msgSuccessDelete, rowID)
	if err != nil {
		log.Print(err)
		message = fmt.Sprintf(msgErrorDelete, rowID)
	}
	return message
}
