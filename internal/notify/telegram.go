// Package notify provides tools for notifing external users.
package notify

import (
	"context"
	"fmt"
	"log"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/tetafro/feed-bot/internal/feed"
)

// Maximum number of concurrent request to telegram api.
const concurrencyLevel = 10

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
		return nil, fmt.Errorf("init telegram API: %w", err)
	}

	bot := &TelegramNotifier{
		api:     api,
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

// Notify notifies all connected client about new event.
func (t *TelegramNotifier) Notify(ctx context.Context, item feed.Item) error {
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

	// TODO: Return real error
	return nil
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
