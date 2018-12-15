package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

// Bot is a telegram bot, that handles two commands: start and stop.
// Starts commands makes bot send greeting message to user every 5 seconds.
// Stop command stops sending messages.
type Bot struct {
	api   *tg.BotAPI
	chats map[int64]*tg.Chat
	mx    sync.RWMutex

	stop chan struct{}
	wg   sync.WaitGroup
}

// NewBot creates new bot.
func NewBot(token string) (*Bot, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "authorization")
	}

	bot := &Bot{
		api:   api,
		chats: make(map[int64]*tg.Chat),
		stop:  make(chan struct{}),
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

	b.wg.Add(1)
	go func() {
		defer b.wg.Done()
		for {
			select {
			case upd := <-updates:
				if err := b.handleCmd(upd); err != nil {
					log.Printf("Failed to handle update: %v", err)
				}
			case <-b.stop:
				return
			}
		}
	}()

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		for range ticker.C {
			b.sendGreetings()
		}
	}()

	return nil
}

// Stop gracefully stops the bot.
func (b *Bot) Stop() {
	close(b.stop)
	b.wg.Wait()
}

func (b *Bot) handleCmd(upd tg.Update) error {
	if upd.Message == nil || !upd.Message.IsCommand() {
		return errors.New("not a command")
	}

	log.Printf("[%s] %s", upd.Message.From.UserName, upd.Message.Text)

	id := upd.Message.Chat.ID
	msg := tg.NewMessage(id, "")
	switch upd.Message.Command() {
	case "start":
		b.addChat(id, upd.Message.Chat)
		msg.Text = "Started"
	case "stop":
		b.removeChat(id)
		msg.Text = "Stopped"
	default:
		msg.Text = "Unknown command"
	}

	_, err := b.api.Send(msg)
	return err
}

func (b *Bot) addChat(id int64, chat *tg.Chat) {
	b.mx.Lock()
	defer b.mx.Unlock()

	b.chats[id] = chat
}

func (b *Bot) removeChat(id int64) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.chats, id)
}

func (b *Bot) sendMessage(id int64, text string) error {
	msg := tg.NewMessage(id, text)
	_, err := b.api.Send(msg)
	if err != nil {
		return errors.Wrap(err, "send message")
	}
	return nil
}

func (b *Bot) sendGreetings() {
	b.mx.RLock()
	defer b.mx.RUnlock()
	for id, chat := range b.chats {
		go func(id int64, chat *tg.Chat) {
			text := fmt.Sprintf(
				"Hello, %s. UTC time is %s.",
				chat.FirstName,
				time.Now().UTC().Format("03:04"),
			)
			if err := b.sendMessage(id, text); err != nil {
				log.Printf("Failed to send message to %d chat: %v", id, err)
			}
		}(id, chat)
	}
}
