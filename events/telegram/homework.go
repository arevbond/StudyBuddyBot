package telegram

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
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

func (p *Processor) AddHomework(text string, userWithChat UserWithChat) string {
	if strings.HasPrefix(text, "/") {
		stateHomework[userWithChat] = newHomework("", "")
		return "Введите название предмета"
	} else if hm, ok := stateHomework[userWithChat]; ok && hm.subject == "" {
		hm.subject = text
		return "Введите задание"
	} else if hm, ok = stateHomework[userWithChat]; ok && hm.Task == "" {
		hm.Task = text
		message := fmt.Sprintf("ДЗ: %s - %s успешно добавлено", hm.subject, hm.Task)
		err := p.storage.AddHomework(context.Background(), userWithChat.ChatID, hm.subject, hm.Task)
		if err != nil {
			message = "Не удалось добавить задание"
			log.Printf("can't add homework: %v", err)
		}
		delete(stateHomework, userWithChat)
		return message
	}
	return "Что-то пошло не так"
}

func (p *Processor) GetHomework(text string, chatID int) string {
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

func (p *Processor) DeleteHomework(rowID int) string {
	err := p.storage.DeleteHomework(context.Background(), rowID)
	message := fmt.Sprintf("Запись №%d успешно удалена", rowID)
	if err != nil {
		log.Print(err)
		message = fmt.Sprintf("Не удалось удалить запись №%d", rowID)
	}
	return message
}
