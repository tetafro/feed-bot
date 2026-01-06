// Telegram bot that reads RSS feeds and sends them to users.
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/sirupsen/logrus"
)

func main() {
	os.Exit(run())
}

func run() int {
	configFile := flag.String("config", "./config.yaml", "path to config file")
	flag.Parse()

	ctx, cancel := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer cancel()

	log := logrus.New()

	conf, err := ReadConfig(*configFile)
	if err != nil {
		log.Errorf("Read config: %v", err)
		return 1
	}

	level := logrus.InfoLevel
	if conf.Debug {
		level = logrus.DebugLevel
	}
	log.SetLevel(level)

	fs, err := NewFileStorage(conf.DataFile)
	if err != nil {
		log.Errorf("Init state storage: %v", err)
		return 1
	}

	tg, err := NewTelegramNotifier(conf.TelegramToken, conf.TelegramChat, log)
	if err != nil {
		log.Errorf("Init telegram notifier: %v", err)
		return 1
	}

	feeds := make([]Feed, len(conf.Feeds))
	for i, url := range conf.Feeds {
		feeds[i] = NewRSSFeed(fs, url)
	}

	log.Info("Starting...")
	NewBot(tg, feeds, conf.UpdateInterval, log).Run(ctx)

	log.Info("Shutdown")
	return 0
}
