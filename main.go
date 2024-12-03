package main

import (
	"io"
	"log"
	"log/slog"
	"os"
	tgClient "tg_ics_useful_bot/clients/telegram"
	"tg_ics_useful_bot/config"
	event_consumer "tg_ics_useful_bot/consumer/event-consumer"
	"tg_ics_useful_bot/events/telegram"
	"tg_ics_useful_bot/storage/cache"
	"tg_ics_useful_bot/storage/postgres"
	"time"

	"github.com/lmittmann/tint"
)

const (
	tgBotHost         = "api.telegram.org"
	storageSQLitePath = "data/sqlite/storage.db"
	batchSize         = 100
)

func main() {
	cfg := config.New()

	logger := setupLogger(cfg)

	s, err := postgres.New(cfg, logger)
	if err != nil {
		logger.Error("can't find storage", slog.Any("err", err))
		os.Exit(1)
	}

	telegramToken := getTelegramToken(cfg)

	eventsProcessor := telegram.New(
		tgClient.New(tgBotHost, telegramToken, cfg.AdminsID),
		s,
		cache.NewUserCache(),
		logger,
	)

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize, logger)

	logger.Info("bot started")

	if err = consumer.Start(); err != nil {
		logger.Error("bot is stopped", slog.Any("err", err))
		os.Exit(1)
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

func setupLogger(cfg *config.Config) *slog.Logger {
	var logger *slog.Logger
	switch cfg.Env {
	case "local":
		f, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		logger = slog.New(tint.NewHandler(io.MultiWriter(f, os.Stdout), &tint.Options{
			Level:      slog.LevelDebug,
			TimeFormat: time.RFC822}))

	case "prod":
		f, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		logger = slog.New(slog.NewJSONHandler(f, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return logger
}
