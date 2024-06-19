package google_calendar

import (
	"context"
	"fmt"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log/slog"
	"os"
	"tg_ics_useful_bot/lib/e"
	"time"
)

const (
	KeyFile       = "credentials.json"
	ScopeCalendar = "https://www.googleapis.com/auth/calendar"
	ScopeEvents   = "https://www.googleapis.com/auth/calendar.events"
)

func Lessons(calendarID string, logger *slog.Logger) (map[time.Weekday][]Lesson, error) {
	lessons := make(map[time.Weekday][]Lesson)
	events, err := allEvents(calendarID, logger)
	if err != nil {
		return nil, err
	}
	items := events.Items
	for _, item := range items {
		if item.Summary != "" && item.Start.DateTime != "" {
			l := rewLesson(item.Summary, item.Start.DateTime, logger)
			lessons[l.DateTime.Weekday()] = append([]Lesson{l}, lessons[l.DateTime.Weekday()]...)
		}
	}
	return lessons, nil
}

type Lesson struct {
	Name     string
	DateTime time.Time
}

func rewLesson(name string, stringTime string, logger *slog.Logger) Lesson {
	t, err := time.Parse(time.RFC3339, stringTime)
	if err != nil {
		logger.Error("can't convert time", slog.Any("error", err), slog.String("string time", stringTime))
	}
	return Lesson{name, t}
}

func allEvents(calendarID string, logger *slog.Logger) (*calendar.Events, error) {
	srv := service(logger)
	events, err := srv.Events.List(calendarID).Do()
	if err != nil {
		return nil, e.Wrap(fmt.Sprintf("can't get events from calendar_id %s", calendarID), err)
	}
	return events, nil
}

func service(logger *slog.Logger) *calendar.Service {
	ctx := context.Background()

	data, err := os.ReadFile("clients/google-calendar/" + KeyFile)
	if err != nil {
		logger.Error("can't read credentials from file", slog.String("file", KeyFile), slog.Any("error", err))
		os.Exit(1)
	}
	creds, err := google.CredentialsFromJSON(ctx, data, ScopeCalendar, ScopeEvents)
	if err != nil {
		logger.Error("can't get credentials from json", slog.Any("error", err))
		os.Exit(1)
	}
	srv, err := calendar.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		logger.Error("Unable to retrieve Calendar client", slog.Any("error", err))
		os.Exit(1)
	}
	return srv
}
