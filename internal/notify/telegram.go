// Package notify provides tools for notifing external users.
package notify

import (
	"context"
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/tetafro/feed-bot/internal/feed"
)

// API describes interface for working with remote API.
type API interface {
	Send(tg.Chattable) (tg.Message, error)
}

// TelegramNotifier uses Telegram API to sends notifications as Telegram
// messages to a public chat.
type TelegramNotifier struct {
	api  API
	chat string
}

// NewTelegramNotifier creates a new bot.
func NewTelegramNotifier(token, chat string) (*TelegramNotifier, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init telegram API: %w", err)
	}

	bot := &TelegramNotifier{
		api:  api,
		chat: "@" + chat,
	}
	return bot, nil
}

// Notify sends a message to the Telegram chat.
func (t *TelegramNotifier) Notify(_ context.Context, item feed.Item) error {
	msg := tg.NewMessageToChannel(t.chat, item.Link)
	_, err := t.api.Send(msg)
	if err != nil {
		return fmt.Errorf("send api request: %w", err)
	}
	return nil
}
