package main

import (
	"context"
	"errors"
	"io"
	"testing"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestTelegramNotifier_Notify(t *testing.T) {
	item := Item{
		Link: "http://example.com/content/",
	}
	log := logrus.New()
	log.Out = io.Discard

	t.Run("successful send", func(t *testing.T) {
		api := &testTgAPI{}
		tn := &TelegramNotifier{api: api, chat: "@chat_name", log: log}

		err := tn.Notify(context.Background(), item)
		assert.NoError(t, err)
		assert.Equal(t, "http://example.com/content/", api.sent)
	})

	t.Run("error from api", func(t *testing.T) {
		api := &testTgAPI{err: errors.New("internal error")}
		tn := &TelegramNotifier{api: api, chat: "@chat_name", log: log}

		err := tn.Notify(context.Background(), item)
		assert.EqualError(t, err, "send api request: internal error")
	})
}

type testTgAPI struct {
	sent string
	err  error
}

func (t *testTgAPI) Send(msg tg.Chattable) (tg.Message, error) {
	switch m := msg.(type) {
	case tg.MessageConfig:
		t.sent = m.Text
	default:
		t.err = errors.New("unknown message type")
	}
	return tg.Message{}, t.err
}
