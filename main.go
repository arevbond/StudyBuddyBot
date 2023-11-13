package main

import (
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	tgClient "tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/config"
	"tg_ics_useful_bot/consumer/event-consumer"
	"tg_ics_useful_bot/events/telegram"
	"tg_ics_useful_bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	storageSQLitePath = "data/sqlite/storage.db"
	batchSize         = 100
)

// init is invoked before main()
func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func main() {
	conf := config.New()

	s, err := sqlite.New(storageSQLitePath)
	if err != nil {
		log.Fatalf("can't connect to storage:", err)
	}

	err = s.Init(context.TODO())
	if err != nil {
		log.Fatalf("can't init storage", err)
	}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, conf.TelegramToken, conf.AdminsID),
		s,
	)

	log.Print("[INFO] service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err = consumer.Start(); err != nil {
		log.Fatal("[ERROR] service is stopped", err)
	}
}

func mustToken() string {
	token := flag.String(
		"tg-bot-token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()

	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *token
}
