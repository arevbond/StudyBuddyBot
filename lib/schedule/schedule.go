package schedule

import (
	"fmt"
	google_calendar "tg_ics_useful_bot/clients/google-calendar"
	"time"
)

func ScheduleCmd(calendarID string) (string, error) {
	result := ""
	dayToLessons, err := google_calendar.Lessons(calendarID)
	if err != nil {
		return "", err
	}
	for day := time.Sunday; day <= time.Saturday; day++ {
		lessons := dayToLessons[day]
		if len(lessons) > 0 {
			result += "\n"
			switch day {
			case time.Monday:
				result += "*Понедельник*"
			case time.Tuesday:
				result += "*Вторник*"
			case time.Wednesday:
				result += "*Среда*"
			case time.Thursday:
				result += "*Четверг*"
			case time.Friday:
				result += "*Пятница*"
			case time.Saturday:
				result += "*Суббота*"
			}
			result += "\n"
			for _, l := range lessons {
				result += fmt.Sprintf("%s - %s\n", l.Name, l.DateTime.Format("15:04"))
			}
		}
	}
	return result, nil
}

func ScheduleByDay(day time.Weekday, calendarID string) (string, error) {
	dayToLessons, err := google_calendar.Lessons(calendarID)
	if err != nil {
		return "", err
	}
	result := ""
	lessons := dayToLessons[day]

	if len(lessons) == 0 {
		return "", nil
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
	return result, nil
}
