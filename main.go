package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	log.Print("Starting...")

	cfg, err := readConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	api, err := tg.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to init telegram API: %v", err)
	}

	bot := NewBot(
		api,
		NewFeed(NewXKCDFetcher(), cfg.UpdateInterval),
		NewFeed(NewCommitStripFetcher(), cfg.UpdateInterval),
		NewFeed(NewExplosmFetcher(), cfg.UpdateInterval),
	)

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
