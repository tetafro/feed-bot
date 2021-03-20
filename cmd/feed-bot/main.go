package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"

	"github.com/tetafro/feed-bot/internal/bot"
	"github.com/tetafro/feed-bot/internal/feed"
	"github.com/tetafro/feed-bot/internal/storage"
)

var configFile = flag.String("f", "./config.json", "path to config file")

func main() {
	log.Print("Starting...")
	if err := run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	log.Print("Shutdown")
}

func run() error {
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := bot.ReadConfig(*configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	api, err := tg.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return errors.Wrap(err, "failed to init telegram API")
	}

	fs, err := storage.NewFileStorage(cfg.DataFile)
	if err != nil {
		return errors.Wrap(err, "failed to init file storage")
	}

	feeds := make([]*feed.Feed, len(cfg.Feeds))
	for i, url := range cfg.Feeds {
		feeds[i] = feed.NewFeed(fs, url, feed.NewRSSFetcher(), cfg.UpdateInterval)
	}

	bot, err := bot.NewBot(api, fs, feeds)
	if err != nil {
		return errors.Wrap(err, "failed to init bot")
	}

	if err := bot.Run(ctx); err != nil {
		return errors.Wrap(err, "failed to start bot")
	}
	return nil
}
