package main

import (
	"flag"
	"log"
	tgClient "tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/config"
	"tg_ics_useful_bot/consumer/event-consumer"
	"tg_ics_useful_bot/events/telegram"
	"tg_ics_useful_bot/storage/postgres"
)

const (
	tgBotHost         = "api.telegram.org"
	storageSQLitePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	cfg := config.New()

	s, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("can't connect to storage:", err)
	}

	//err = s.Init(context.TODO())
	//if err != nil {
	//	log.Fatalf("can't init storage", err)
	//}

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, cfg.TelegramToken, cfg.AdminsID),
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
