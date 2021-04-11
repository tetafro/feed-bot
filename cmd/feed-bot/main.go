package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/pkg/errors"

	"github.com/tetafro/feed-bot/internal/bot"
	"github.com/tetafro/feed-bot/internal/feed"
	"github.com/tetafro/feed-bot/internal/notify"
	"github.com/tetafro/feed-bot/internal/storage"
)

var configFile = flag.String("f", "./config.yaml", "path to config file")

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

	conf, err := bot.ReadConfig(*configFile)
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	fs, err := storage.NewFileStorage(conf.DataFile)
	if err != nil {
		return errors.Wrap(err, "failed to init state storage")
	}

	var wg sync.WaitGroup
	var notifier bot.Notifier
	if conf.LogNotifier {
		notifier = notify.NewLogNotifier()
	} else {
		tg, err := notify.NewTelegramNotifier(conf.TelegramToken, fs)
		if err != nil {
			return errors.Wrap(err, "failed to init telegram notifier")
		}
		wg.Add(1)
		go func() {
			tg.ListenCommands(ctx)
			wg.Done()
		}()
		notifier = tg
	}

	feeds := make([]bot.Feed, len(conf.Feeds))
	for i, url := range conf.Feeds {
		feeds[i] = feed.NewRSSFeed(fs, url)
	}

	bot.New(notifier, feeds, conf.UpdateInterval).Run(ctx)

	wg.Wait()
	return nil
}
