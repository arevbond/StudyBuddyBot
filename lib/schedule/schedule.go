package schedule

import (
	"fmt"
	"log/slog"
	google_calendar "tg_ics_useful_bot/clients/google-calendar"
	"time"
)

func ScheduleCmd(calendarID string, logger *slog.Logger) (string, error) {
	result := ""
	dayToLessons, err := google_calendar.Lessons(calendarID, logger)
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
