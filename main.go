package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("Starting...")

	cfg := mustConfig()

	bot, err := NewBot(cfg.TelegramToken)
	if err != nil {
		log.Fatalf("Failed to init bot: %v", err)
	}

	if err := bot.Start(); err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	log.Print("Ready to work")

	handleSignals()
	bot.Stop()
	log.Print("Shutdown")
}

func handleSignals() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stop)

	<-stop
}
