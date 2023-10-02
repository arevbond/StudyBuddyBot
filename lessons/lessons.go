package lessons

import "time"

const (
	msgSunday = "Воскресенье же....🙈"
)

type lesson struct {
	Name    string
	Teacher string
	Time    string
}

func TomorrowLessons() string {
	result := "Расписание на завтра:\n\n"
	result += stringTomorrowLessons(time.Now().Weekday())
	return result
}

func AllLessons() string {
	result := "Расписание на неделю:\n\n"
	result += stringAllLessons()
	return result
}

func LessonsToday() string {
	result := "Расписание на сегодня:\n\n"
	result += stringLessonsByDay(time.Now().Weekday())
	return result
}

func Lessons() map[time.Weekday][]lesson {
	lessons := make(map[time.Weekday][]lesson)
	lessons[time.Monday] = []lesson{
		{Name: "БЖД", Teacher: "Дубова", Time: "9:45"},
		{Name: "БЖД", Teacher: "Хорошилова", Time: "11:30"},
	}
	lessons[time.Tuesday] = []lesson{
		{Name: "Аналитика больших объёмов данных", Teacher: "Прохоров Кирюха", Time: "8:00"},
		{Name: "Аналитика больших объёмов данных", Teacher: "Прохоров Кирюха", Time: "9:45"},
	}
	lessons[time.Wednesday] = []lesson{
		{Name: "Психология личности и ее саморазвития", Teacher: "Велимедова", Time: "--:--"},
	}

	lessons[time.Thursday] = []lesson{
		{Name: "ОВП", Teacher: "Зайцев", Time: "8:00"},
	}
	lessons[time.Friday] = []lesson{
		{Name: "Основы права", Teacher: "Саприн И.Г.", Time: "9:45"},
		{Name: "ООП", Teacher: "Коровченко И.С", Time: "11:30"},
	}
	lessons[time.Saturday] = []lesson{
		{Name: "Защита информации", Teacher: "Овцинникова Т.М.", Time: "9:45"},
		{Name: "Основы теории передачи информации", Teacher: "Гутерман Н.Е.", Time: "13:25"},
	}
	return lessons
}

func stringLessonsByDay(day time.Weekday) string {
	ls := Lessons()
	result := ""
	today := ls[day]
	if today == nil {
		result += msgSunday
	}
	for _, l := range today {
		result += "• " + "Время: " + l.Time + ". " + l.Name + ". " + l.Teacher + "\n"
	}
	return result
}

func stringAllLessons() string {
	result := ""
	for d := time.Monday; d <= time.Saturday; d++ {
		switch d {
		case time.Monday:
			result += "Понедельник:\n"
		case time.Tuesday:
			result += "Вторник:\n"
		case time.Wednesday:
			result += "Среда:\n"
		case time.Thursday:
			result += "Четверг:\n"
		case time.Friday:
			result += "Пятница:\n"
		case time.Saturday:
			result += "Суббота:\n"
		}
		result += stringLessonsByDay(d)
		result += "\n"
	}
	return result
}

func stringTomorrowLessons(day time.Weekday) string {
	day += 1
	if day == 7 {
		day = 0
	}
	return stringLessonsByDay(day)
}
