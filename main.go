// Telegram bot that reads RSS feeds and sends them to users.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Print("Starting...")
	if err := run(); err != nil {
		log.Fatalf("ERROR: %v", err)
	}
	log.Print("Shutdown")
}

func run() error {
	configFile := flag.String("config", "./config.yaml", "path to config file")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	conf, err := ReadConfig(*configFile)
	if err != nil {
		return fmt.Errorf("read config: %w", err)
	}

	fs, err := NewFileStorage(conf.DataFile)
	if err != nil {
		return fmt.Errorf("init state storage: %w", err)
	}

	var notifier Notifier
	if conf.LogNotifier {
		notifier = NewLogNotifier()
	} else {
		tg, err := NewTelegramNotifier(conf.TelegramToken, conf.TelegramChat)
		if err != nil {
			return fmt.Errorf("init telegram notifier: %w", err)
		}
		notifier = tg
	}

	feeds := make([]Feed, len(conf.Feeds))
	for i, url := range conf.Feeds {
		feeds[i] = NewRSSFeed(fs, url)
	}

	NewBot(notifier, feeds, conf.UpdateInterval).Run(ctx)

	return nil
}
