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

type Lesson struct {
	Name     string
	DateTime time.Time
}

func NewLesson(name string, stringTime string) Lesson {
	t, err := time.Parse(time.RFC3339, stringTime)
	if err != nil {
		log.Printf("can't convert time from string %v: %v", stringTime, err)
	}
	return Lesson{name, t}
}

type Manager struct {
	CalendarID string
	srv        *calendar.Service
}

func NewManager(calendarID string) Manager {
	return Manager{calendarID, service()}
}

func (m Manager) Lessons() map[time.Weekday]Lesson {
	lessons := make(map[time.Weekday]Lesson)
	events := m.allEvents()
	items := events.Items
	for _, item := range items {
		l := NewLesson(item.Summary, item.Start.DateTime)
		lessons[l.DateTime.Weekday()] = l
	}
	return lessons
}

func (m Manager) allEvents() *calendar.Events {
	events, err := m.srv.Events.List(m.CalendarID).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve next ten of the user's events: %v", err)
	}
	return events
}

func service() *calendar.Service {
	ctx := context.Background()

	data, err := os.ReadFile(KeyFile)
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
