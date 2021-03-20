// Package bot provides main application entity. It is responsible for wiring
// together all components: data feeds, chats, state storage.
package bot

import (
	"context"
	"log"
	"sync"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"

	"github.com/tetafro/feed-bot/internal/feed"
)

// Storage describes persistent datastorage.
type Storage interface {
	GetChats() (ids []int64, err error)
	SaveChats(ids []int64) error
}

// API describes interface for working with remote API.
type API interface {
	GetUpdatesChan(tg.UpdateConfig) (tg.UpdatesChannel, error)
	Send(tg.Chattable) (tg.Message, error)
}

// Bot is a telegram bot, that handles two commands: start and stop.
// Starts commands makes bot send feed updates to user.
// Stop command stops sending messages.
type Bot struct {
	// API for sending messages
	api API

	// Set of RSS feeds
	feeds []*feed.Feed

	// List of connected users
	chats map[int64]struct{}
	mx    *sync.Mutex

	// Storage for current chats data
	storage Storage
}

// NewBot creates new bot.
func NewBot(api API, st Storage, feeds []*feed.Feed) (*Bot, error) {
	bot := &Bot{
		api:     api,
		feeds:   feeds,
		chats:   map[int64]struct{}{},
		mx:      &sync.Mutex{},
		storage: st,
	}

	chatIDs, err := st.GetChats()
	if err != nil {
		return nil, errors.Wrap(err, "get chats data")
	}
	if len(chatIDs) == 0 {
		log.Print("No chats data")
		return bot, nil
	}

	log.Printf("Currently connected users: %d", len(chatIDs))

	for _, id := range chatIDs {
		bot.chats[id] = struct{}{}
	}
	return bot, nil
}

// Run starts listening for updates.
func (b *Bot) Run(ctx context.Context) error {
	updates, err := b.api.GetUpdatesChan(tg.UpdateConfig{
		Offset:  0,
		Limit:   0,
		Timeout: 30,
	})
	if err != nil {
		return errors.Wrap(err, "get updates channel")
	}

	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		b.listenCommands(ctx, updates)
		wg.Done()
	}()
	go func() {
		b.processFeed()
		wg.Done()
	}()

	wg.Add(len(b.feeds))
	for _, f := range b.feeds {
		f := f
		go func() {
			f.Run(ctx)
			wg.Done()
		}()
	}

	wg.Wait()
	return nil
}

func (b *Bot) listenCommands(ctx context.Context, updates tg.UpdatesChannel) {
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
		case <-ctx.Done():
			return
		}
	}
}

func (b *Bot) processFeed() {
	var wg sync.WaitGroup
	for item := range b.feed() {
		for chat := range b.chats {
			chat := chat
			wg.Add(1)
			go func() {
				b.send(chat, item)
				wg.Done()
			}()
		}
	}
	wg.Wait()
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

	if err := b.storage.SaveChats(mapToSlice(b.chats)); err != nil {
		log.Printf("Failed to save chats data: %v", err)
	}
}

func (b *Bot) removeChat(id int64) {
	b.mx.Lock()
	defer b.mx.Unlock()

	delete(b.chats, id)

	if err := b.storage.SaveChats(mapToSlice(b.chats)); err != nil {
		log.Printf("Failed to save chats data: %v", err)
	}
}

// feed merges all feeds into one channel.
func (b *Bot) feed() <-chan feed.Item {
	wg := sync.WaitGroup{}
	wg.Add(len(b.feeds))

	all := make(chan feed.Item)
	for _, f := range b.feeds {
		go func(ch <-chan feed.Item) {
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

func (b *Bot) send(chat int64, item feed.Item) {
	msg := tg.NewPhotoShare(chat, item.Image)
	msg.Caption = item.Title
	_, err := b.api.Send(msg)
	if err != nil {
		log.Printf("Failed to send message: %v", err)
	}
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
