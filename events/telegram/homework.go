package telegram

import (
	"fmt"
	"strconv"
	"strings"
)

type Homework struct {
	subject string
	data    string
}

func newHomework(subject, data string) Homework {
	return Homework{subject: subject, data: data}
}

var chatToHomewords = make(map[int][]Homework)

func (p *Processor) AddHomework(text string, chatID int) string {
	subject := ""
	data := ""
	for _, str := range strings.Split(text, " ")[1:] {
		if strings.HasPrefix(str, "#") {
			subject = str[1:]
		} else if str != "" {
			data += str + " "
		}
	}
	if subject == "" {
		return msgHomeworkWithoutSubject
	} else if data == "" {
		return msgHomeworkWithoutData
	}
	homework := newHomework(subject, data)
	chatToHomewords[chatID] = append(chatToHomewords[chatID], homework)
	return fmt.Sprintf(msgHomeworkSuccessAdded, homework.subject, homework.data)
}

func (p *Processor) GetHomework(text string, chatID int) string {
	val := ""
	for _, s := range strings.Split(text, " ")[1:] {
		if s != "" {
			val = s
			break
		}
	}
	var isDigit bool
	resultHomeworks := make([]Homework, 0)
	homeworks := chatToHomewords[chatID]
	if strings.HasPrefix(val, "#") {
		for _, hm := range homeworks {
			if hm.subject == val[1:] {
				resultHomeworks = append(resultHomeworks, hm)
			}
		}
		isDigit = false
	} else if num, err := strconv.Atoi(val); err == nil {
		for i := len(homeworks) - 1; i >= 0 && num > 0; i-- {
			resultHomeworks = append(resultHomeworks, homeworks[i])
			num--
		}
		isDigit = true
	}

	message := ""
	if isDigit {
		num, _ := strconv.Atoi(val)
		if num > len(resultHomeworks) {
			num = len(resultHomeworks)
		}
		message += fmt.Sprintf("Последние %d домашних заданий:\n", num)
	} else if val != "" {
		message += fmt.Sprintf("Домашнее задания по предмету %s:\n", val)
	} else {
		for i := len(homeworks) - 1; i >= 0; i-- {
			resultHomeworks = append(resultHomeworks, homeworks[i])
		}
		message += fmt.Sprintf("Всё домашнее задание:\n")
	}
	for _, h := range resultHomeworks {
		message += fmt.Sprintf(" • %s - %s\n", h.subject, h.data)
	}
	return message
}
