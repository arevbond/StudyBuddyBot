package event_consumer

import (
	"log/slog"
	"tg_ics_useful_bot/events"
	"time"
)

type Consumer struct {
	logger    *slog.Logger
	fetcher   events.Fetcher
	processor events.Processor
	batchSize int
}

func New(fetcher events.Fetcher, processor events.Processor, batchSize int, logger *slog.Logger) Consumer {
	return Consumer{
		logger:    logger,
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

func (c Consumer) Start() error {
	for {
		gotEvents, err := c.fetcher.Fetch(c.batchSize)
		if err != nil {
			c.logger.Error("can't fetch events", slog.Any("error", err))

			continue
		}

		if len(gotEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err = c.handleEvents(gotEvents); err != nil {
			c.logger.Error("consumer can't handle events", slog.Any("error", err))

			continue
		}
	}
}

func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		if err := c.processor.Process(event); err != nil {
			c.logger.Error("can't handle events", slog.Any("error", err))
			continue
		}

	}
	return nil
}
