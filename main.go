package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

var configFile = flag.String("f", "./config.json", "path to config file")

func main() {
	flag.Parse()

	log.Print("Starting...")

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	cfg, err := readConfig(*configFile)
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	api, err := tg.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to init telegram API: %v", err)
	}

	fs, err := NewFileStorage(cfg.DataFile)
	if err != nil {
		log.Fatalf("Failed to init file storage: %v", err)
	}

	feeds := make([]*Feed, len(cfg.Feeds))
	i := 0
	for _, url := range cfg.Feeds {
		feeds[i] = NewFeed(fs, url, NewRSSFetcher(), cfg.UpdateInterval)
		i++
	}

	bot, err := NewBot(api, fs, feeds...)
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}

	if err := bot.Run(ctx); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	log.Print("Shutdown")
}
