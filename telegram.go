package main

import (
	"context"
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/sirupsen/logrus"
)

// API describes interface for working with remote API.
type API interface {
	Send(tg.Chattable) (tg.Message, error)
}

// TelegramNotifier uses Telegram API to sends notifications as Telegram
// messages to a channel.
type TelegramNotifier struct {
	api  API
	chat string
	log  *logrus.Logger
}

// NewTelegramNotifier creates a new telegram client.
func NewTelegramNotifier(token, chat string, log *logrus.Logger) (*TelegramNotifier, error) {
	api, err := tg.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("init telegram API: %w", err)
	}

	bot := &TelegramNotifier{
		api:  api,
		chat: "@" + chat,
		log:  log,
	}
	return bot, nil
}

// Notify sends a message to a Telegram channel.
func (t *TelegramNotifier) Notify(_ context.Context, item Item) error {
	msg := tg.NewMessageToChannel(t.chat, item.Link)
	_, err := t.api.Send(msg)
	if err != nil {
		return fmt.Errorf("send api request: %w", err)
	}
	return nil
}

// PrintNotifier is a notifier for debugging mode. It prints items without
// sending them anywhere.
type PrintNotifier struct {
	log *logrus.Logger
}

// NewPrintNotifier creates a new PrintNotifier.
func NewPrintNotifier(log *logrus.Logger) *PrintNotifier {
	return &PrintNotifier{log: log}
}

// Notify prints the item.
//
//nolint:unparam
func (n *PrintNotifier) Notify(_ context.Context, item Item) error {
	n.log.Debugf("New item: %s", item)
	return nil
}
