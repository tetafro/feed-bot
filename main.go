package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("Starting...")

	cfg, err := readConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	bot, err := NewBot(
		cfg.TelegramToken,
		NewFeed(NewXKCDFetcher(), cfg.UpdateInterval),
		NewFeed(NewCommitStripFetcher(), cfg.UpdateInterval),
		NewFeed(NewExplosmFetcher(), cfg.UpdateInterval),
	)
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
