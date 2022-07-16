package notify

import (
	"bytes"
	"context"
	"log"
	"sort"
	"sync"
	"testing"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"

	"github.com/tetafro/feed-bot/internal/feed"
)

func TestTelegramNotifier_ListenCommands(t *testing.T) {
	t.Run("add chat", func(t *testing.T) {
		defer log.SetOutput(log.Writer())
		defer log.SetFlags(log.Flags())

		var buf bytes.Buffer
		log.SetOutput(&buf)
		log.SetFlags(0)

		ctx, cancel := context.WithCancel(context.Background())

		cmd := make(chan tg.Update)
		api := &testTgAPI{}
		tn := &TelegramNotifier{
			api:     api,
			cmd:     cmd,
			chats:   map[int64]struct{}{},
			mx:      &sync.Mutex{},
			storage: &testStorage{},
		}

		go func() {
			ee := []tg.MessageEntity{{
				Type:   "bot_command",
				Length: 6,
			}}
			cmd <- tg.Update{
				Message: &tg.Message{
					Chat:     &tg.Chat{ID: 100},
					From:     &tg.User{UserName: "user"},
					Text:     "/start",
					Entities: &ee,
				},
			}
			cancel()
		}()

		tn.ListenCommands(ctx)

		assert.Equal(t, "[user] /start\n", buf.String())
		assert.Equal(t, map[int64]struct{}{100: {}}, tn.chats)
		assert.Equal(t, "Started", api.sent)
	})

	t.Run("remove chat", func(t *testing.T) {
		defer log.SetOutput(log.Writer())
		defer log.SetFlags(log.Flags())

		var buf bytes.Buffer
		log.SetOutput(&buf)
		log.SetFlags(0)

		ctx, cancel := context.WithCancel(context.Background())

		cmd := make(chan tg.Update)
		api := &testTgAPI{}
		tn := &TelegramNotifier{
			api:     api,
			cmd:     cmd,
			chats:   map[int64]struct{}{100: {}},
			mx:      &sync.Mutex{},
			storage: &testStorage{},
		}

		go func() {
			ee := []tg.MessageEntity{{
				Type:   "bot_command",
				Length: 5,
			}}
			cmd <- tg.Update{
				Message: &tg.Message{
					Chat:     &tg.Chat{ID: 100},
					From:     &tg.User{UserName: "user"},
					Text:     "/stop",
					Entities: &ee,
				},
			}
			cancel()
		}()

		tn.ListenCommands(ctx)

		assert.Equal(t, "[user] /stop\n", buf.String())
		assert.Equal(t, map[int64]struct{}{}, tn.chats)
		assert.Equal(t, "Stopped", api.sent)
	})

	t.Run("unknown command", func(t *testing.T) {
		defer log.SetOutput(log.Writer())
		defer log.SetFlags(log.Flags())

		var buf bytes.Buffer
		log.SetOutput(&buf)
		log.SetFlags(0)

		ctx, cancel := context.WithCancel(context.Background())

		cmd := make(chan tg.Update)
		tn := &TelegramNotifier{
			cmd:   cmd,
			chats: map[int64]struct{}{},
			mx:    &sync.Mutex{},
		}
		go func() {
			cmd <- tg.Update{}
			cancel()
		}()
		tn.ListenCommands(ctx)

		assert.Contains(t, buf.String(), "Unknown command")
	})
}

func TestTelegramNotifier_Notify(t *testing.T) {
	t.Run("send to all", func(t *testing.T) {
		api := &testTgAPI{}
		tn := &TelegramNotifier{
			api:   api,
			chats: map[int64]struct{}{100: {}},
			mx:    &sync.Mutex{},
		}

		_ = tn.Notify(context.Background(), feed.Item{
			Link: "http://example.com/content/",
		})

		assert.Equal(t, "http://example.com/content/", api.sent)
	})

	t.Run("no active chats", func(t *testing.T) {
		api := &testTgAPI{}
		tn := &TelegramNotifier{
			api:   api,
			chats: map[int64]struct{}{},
			mx:    &sync.Mutex{},
		}

		_ = tn.Notify(context.Background(), feed.Item{
			Link: "http://example.com/content/",
		})

		assert.Equal(t, "", api.sent)
	})
}

func TestMapToSlice(t *testing.T) {
	testCases := []struct {
		name string
		m    map[int64]struct{}
		s    []int64
	}{
		{
			name: "test-1",
			m:    map[int64]struct{}{1: {}, 2: {}, 3: {}},
			s:    []int64{1, 2, 3},
		},
		{
			name: "test-2",
			m:    map[int64]struct{}{1: {}},
			s:    []int64{1},
		},
		{
			name: "test-3",
			m:    map[int64]struct{}{},
			s:    []int64{},
		},
		{
			name: "test-4",
			m:    nil,
			s:    nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			s := mapToSlice(tt.m)
			sort.Slice(s, func(i, j int) bool {
				return s[i] < s[j]
			})
			assert.Equal(t, tt.s, s)
		})
	}
}

type testTgAPI struct {
	sent string
}

func (t *testTgAPI) GetUpdatesChan(tg.UpdateConfig) (tg.UpdatesChannel, error) {
	return nil, nil
}

func (t *testTgAPI) Send(msg tg.Chattable) (tg.Message, error) {
	switch m := msg.(type) {
	case tg.MessageConfig:
		t.sent = m.Text
	case tg.PhotoConfig:
		t.sent = m.Caption + "|" + m.BaseFile.FileID
	}
	return tg.Message{}, nil
}

type testStorage struct{}

func (s *testStorage) GetChats() []int64 {
	return []int64{100}
}

func (s *testStorage) SaveChats(ids []int64) error {
	return nil
}
