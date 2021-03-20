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

	cfg, err := readConfig(*configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	api, err := tg.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		return errors.Wrap(err, "failed to init telegram API")
	}

	fs, err := NewFileStorage(cfg.DataFile)
	if err != nil {
		return errors.Wrap(err, "failed to init file storage")
	}

	feeds := make([]*Feed, len(cfg.Feeds))
	i := 0
	for _, url := range cfg.Feeds {
		feeds[i] = NewFeed(fs, url, NewRSSFetcher(), cfg.UpdateInterval)
		i++
	}

	bot, err := NewBot(api, fs, feeds...)
	if err != nil {
		return errors.Wrap(err, "failed to init bot")
	}

	if err := bot.Run(ctx); err != nil {
		return errors.Wrap(err, "failed to start bot")
	}
	return nil
}
