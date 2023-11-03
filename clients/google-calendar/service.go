package google_calendar

import (
	"context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
	"log"
	"os"
	"time"
)

const (
	KeyFile       = "credentials.json"
	ScopeCalendar = "https://www.googleapis.com/auth/calendar"
	ScopeEvents   = "https://www.googleapis.com/auth/calendar.events"
)

func Lessons(calendarID string) map[time.Weekday][]Lesson {
	lessons := make(map[time.Weekday][]Lesson)
	events := allEvents(calendarID)
	items := events.Items
	for _, item := range items {
		l := rewLesson(item.Summary, item.Start.DateTime)
		lessons[l.DateTime.Weekday()] = append(lessons[l.DateTime.Weekday()], l)
	}
	return lessons
}

type Lesson struct {
	Name     string
	DateTime time.Time
}

func rewLesson(name string, stringTime string) Lesson {
	t, err := time.Parse(time.RFC3339, stringTime)
	if err != nil {
		log.Printf("can't convert time from string %v: %v", stringTime, err)
	}
	return Lesson{name, t}
}

func allEvents(calendarID string) *calendar.Events {
	srv := service()
	events, err := srv.Events.List(calendarID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	return events
}

func service() *calendar.Service {
	ctx := context.Background()

	data, err := os.ReadFile("clients/google-calendar/" + KeyFile)
	if err != nil {
		log.Fatalf("Can't read credentials from file: %s", KeyFile)
	}
	creds, err := google.CredentialsFromJSON(ctx, data, ScopeCalendar, ScopeEvents)
	if err != nil {
		log.Fatal(err)
	}
	srv, err := calendar.NewService(ctx, option.WithCredentials(creds))
	if err != nil {
		log.Fatalf("Unable to retrieve Calendar client: %v", err)
	}
	return srv
}
