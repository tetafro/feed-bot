package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

const configFile = "./config.json"

func main() {
	log.Print("Starting...")

	cfg, err := readConfig(configFile)
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

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	log.Print("Ready to work")

	handleSignals()
	log.Print("Shutdown")
}

func handleSignals() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	sig := <-stop
	log.Printf("Got %s", sig)
}
