// Package notify provides tools for notifing external users.
package notify

import (
	"context"
	"log"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"

	"github.com/tetafro/feed-bot/internal/feed"
)

const (
	// List of available telegram bot commands.
	startCmd = "start"
	stopCmd  = "stop"

	// Maximum number of concurrent request to telegram api.
	concurrencyLevel = 10
)

// API describes interface for working with remote API.
type API interface {
	GetUpdatesChan(tg.UpdateConfig) (tg.UpdatesChannel, error)
	Send(tg.Chattable) (tg.Message, error)
}

// TelegramNotifier is a telegram bot, that handles two commands: start and
// stop.
// Starts commands makes bot send feed updates to user.
// Stop command stops sending messages.
type TelegramNotifier struct {
	// API for sending messages
	api API

	// Incoming commands from telegram
	cmd tg.UpdatesChannel

	// List of connected users
	chats map[int64]struct{}
	mx    *sync.Mutex

	// Storage for current chats data
	storage Storage
}

// NewTelegramNotifier creates new bot.
func NewTelegramNotifier(token string, st Storage) (*TelegramNotifier, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, errors.Wrap(err, "init telegram API")
	}
	cmd, err := api.GetUpdatesChan(tg.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 30,
	})
	if err != nil {
		return nil, errors.Wrap(err, "get commands channel")
	}

	bot := &TelegramNotifier{
		api:     api,
		cmd:     cmd,
		chats:   map[int64]struct{}{},
		mx:      &sync.Mutex{},
		storage: st,
	}

	chats := st.GetChats()
	for _, id := range chats {
		bot.chats[id] = struct{}{}
	}
	log.Printf("Currently connected users: %d", len(chats))
	return bot, nil
}

// ListenCommands starts listening to user commands.
func (t *TelegramNotifier) ListenCommands(ctx context.Context) {
	for {
		select {
		case cmd := <-t.cmd:
			if cmd.Message == nil || !cmd.Message.IsCommand() {
				log.Printf("Unknown command")
				continue
			}
			if err := t.handleCommand(cmd); err != nil {
				log.Printf("Failed to handle command: %v", err)
			}
		case <-ctx.Done():
			return
		}
	}
}

// Notify notifies all connected client about new event.
func (t *TelegramNotifier) Notify(ctx context.Context, item feed.Item) {
	t.mx.Lock()
	chats := mapToSlice(t.chats)
	t.mx.Unlock()

	// Setup semaphore to control concurrency
	sema := make(chan struct{}, concurrencyLevel)
	for i := 0; i < concurrencyLevel; i++ {
		sema <- struct{}{}
	}
	// Send message to each connected client
	for _, c := range chats {
		<-sema
		c := c
		go func() {
			if err := t.send(c, item); err != nil {
				log.Printf("Failed to send message: %v", err)
			}
			sema <- struct{}{}
		}()
	}
	// Wait for all goroutines
	for i := 0; i < concurrencyLevel; i++ {
		<-sema
	}
}

func (t *TelegramNotifier) handleCommand(upd tg.Update) error {
	log.Printf("[%s] %s", upd.Message.From.UserName, upd.Message.Text)
	cid := upd.Message.Chat.ID

	var text string
	switch upd.Message.Command() {
	case startCmd:
		text = "Started"
		if err := t.addChat(cid); err != nil {
			log.Printf("Failed to save chats data: %v", err)
		}
	case stopCmd:
		text = "Stopped"
		if err := t.removeChat(cid); err != nil {
			log.Printf("Failed to save chats data: %v", err)
		}
	default:
		text = "Unknown command"
	}

	_, err := t.api.Send(tg.NewMessage(cid, text))
	return err
}

func (t *TelegramNotifier) addChat(id int64) error {
	t.mx.Lock()
	defer t.mx.Unlock()

	t.chats[id] = struct{}{}

	return t.storage.SaveChats(mapToSlice(t.chats))
}

func (t *TelegramNotifier) removeChat(id int64) error {
	t.mx.Lock()
	defer t.mx.Unlock()

	delete(t.chats, id)

	return t.storage.SaveChats(mapToSlice(t.chats))
}

func (t *TelegramNotifier) send(chat int64, item feed.Item) error {
	msg := tg.NewMessage(chat, item.Link)
	_, err := t.api.Send(msg)
	return err
}

func mapToSlice(m map[int64]struct{}) []int64 {
	if m == nil {
		return nil
	}
	nn := make([]int64, len(m))
	i := 0
	for n := range m {
		nn[i] = n
		i++
	}
	return nn
}
