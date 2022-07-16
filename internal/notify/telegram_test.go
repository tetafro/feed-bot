package notify

import (
	"context"
	"sort"
	"sync"
	"testing"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"

	"github.com/tetafro/feed-bot/internal/feed"
)

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
