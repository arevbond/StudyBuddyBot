package main

import (
	"log"
	tgClient "tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/config"
	"tg_ics_useful_bot/consumer/event-consumer"
	"tg_ics_useful_bot/events/telegram"
	"tg_ics_useful_bot/storage/cache"
	"tg_ics_useful_bot/storage/sqlite"
)

const (
	tgBotHost         = "api.telegram.org"
	storageSQLitePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	cfg := config.New()

	s, err := sqlite.New(storageSQLitePath)
	if err != nil {
		log.Fatal("can't find storage", err)
	}

	telegramToken := getTelegramToken(cfg)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, telegramToken, cfg.AdminsID),
		s,
		cache.NewUserCache(),
	)

	log.Print("[INFO] bot started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err = consumer.Start(); err != nil {
		log.Fatal("[ERROR] bot is stopped", err)
	}
}

func getTelegramToken(cfg *config.Config) string {
	var telegramToken string

	switch cfg.Env {
	case "local":
		telegramToken = cfg.TestTelegramToken
	case "prod":
		telegramToken = cfg.TelegramToken
	}

	return telegramToken
}
