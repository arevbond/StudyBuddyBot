package schedule

import (
	"fmt"
	google_calendar "tg_ics_useful_bot/clients/google-calendar"
	"time"
)

func Schedule(calendarID string) string {
	result := "Расписание:\n\n"
	for day := time.Monday; day <= time.Saturday; day++ {
		cur := scheduleByDay(day, calendarID)
		if cur != "" {
			result += cur + "\n"
		}
	}
	return result
}

func scheduleByDay(day time.Weekday, calendarID string) string {
	dayToLessons := google_calendar.Lessons(calendarID)
	result := ""
	lessons := dayToLessons[day]

	if len(lessons) == 0 {
		return ""
	}

	switch day {
	case time.Monday:
		result += "Понедельник"
	case time.Tuesday:
		result += "Вторник"
	case time.Wednesday:
		result += "Среда"
	case time.Thursday:
		result += "Четверг"
	case time.Friday:
		result += "Пятница"
	case time.Saturday:
		result += "Суббота"
	}
	result += "\n"

	for _, l := range lessons {
		result += fmt.Sprintf("%s - %s\n", l.Name, l.DateTime.Format("15:04"))
	}
	return result
}
