package main

import (
	"log"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// Bot is a telegram bot, that handles two commands: start and stop.
// Starts commands makes bot send greeting message to user every 5 seconds.
// Stop command stops sending messages.
type Bot struct {
	api *tg.BotAPI

	feeds []Feed

	chats map[int64]struct{}
	mx    *sync.Mutex

	stop chan struct{}
	wg   *sync.WaitGroup
}

// NewBot creates new bot.
func NewBot(token string, feeds []Feed) (*Bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "authorization")
	}

	bot := &Bot{
		api:   api,
		feeds:  feeds,
		chats: map[int64]struct{}{},
		mx:    &sync.Mutex{},
		stop:  make(chan struct{}),
		wg:    &sync.WaitGroup{},
	}

	return bot, nil
}

// Start starts listening for updates.
func (b *Bot) Start() error {
	updates, err := b.api.GetUpdatesChan(tg.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 30,
	})
	if err != nil {
		return errors.Wrap(err, "get updates channel")
	}

	b.wg.Add(2)
	go b.listenCommands(updates)
	go b.processFeed()

	for _, f := range b.feeds {
	f.Start()
	}

	return nil
}

// Stop gracefully stops the bot.
func (b *Bot) Stop() {
	for _, f := range b.feeds {
		f.Stop()
	}
	close(b.stop)
	b.wg.Wait()
}

func (b *Bot) listenCommands(updates tg.UpdatesChannel) {
	defer b.wg.Done()
	for {
		select {
		case upd := <-updates:
			if upd.Message == nil || !upd.Message.IsCommand() {
				log.Printf("Unknown update type")
				continue
			}
			if err := b.handleCommand(upd); err != nil {
				log.Printf("Failed to handle update: %v", err)
			}
		case <-b.stop:
			return
		}
	}
}

func (b *Bot) processFeed() {
	defer b.wg.Done()

	for item := range b.feed() {
		for chat := range b.chats {
			msg := tg.NewMessage(chat, item.String())
			_, err := b.api.Send(msg)
			if err != nil {
				log.Printf("Failed to send message: %v", err)
			}
		}
	}
}

func (b *Bot) handleCommand(upd tg.Update) error {
	log.Printf("[%s] %s", upd.Message.From.UserName, upd.Message.Text)
	chatID := upd.Message.Chat.ID

	var text string
	switch upd.Message.Command() {
	case "start":
		text = "Started"
		b.addChat(chatID)
	case "stop":
		text = "Stopped"
		b.removeChat(chatID)
	case "xkcd":
		item, err := b.feeds[0].Last()
		if err != nil {
			text = "Service error, try again later"
			break
		}
		text = item.String()
	default:
		text = "Unknown command"
	}

	msg := tg.NewMessage(chatID, text)
	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) addChat(id int64) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.chats[id] = struct{}{}
}

func (b *Bot) removeChat(id int64) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.chats, id)
}

// feed merges all feeds into one channel.
func (b *Bot) feed() <-chan Item {
	wg := sync.WaitGroup{}
	wg.Add(len(b.feeds))

	all := make(chan Item)
	for _, f := range b.feeds {
		go func(ch <-chan Item) {
			defer wg.Done()
			for item := range ch {
				all <- item
			}
		}(f.Feed())
	}

	go func() {
		wg.Wait()
		close(all)
	}()

	return all
}
