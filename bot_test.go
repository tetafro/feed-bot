package main

import (
	"testing"

	tg "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/stretchr/testify/assert"
)

func TestNewBot(t *testing.T) {
	api := &mockAPI{}
	st := &mockStorage{}
	bot, err := NewBot(api, st)
	assert.NoError(t, err)
	assert.NotNil(t, bot)
}

func TestMapToSlice(t *testing.T) {
	testCases := []struct {
		name string
		m    map[int64]struct{}
		s    []int64
	}{
		{
			name: "test-1",
			m:    map[int64]struct{}{1: struct{}{}, 2: struct{}{}, 3: struct{}{}},
			s:    []int64{1, 2, 3},
		},
		{
			name: "test-2",
			m:    map[int64]struct{}{1: struct{}{}},
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
			assert.Equal(t, tt.s, mapToSlice(tt.m))
		})
	}
}

type mockAPI struct{}

func (m *mockAPI) GetUpdatesChan(tg.UpdateConfig) (tg.UpdatesChannel, error) {
	return nil, nil
}

func (m *mockAPI) Send(tg.Chattable) (tg.Message, error) {
	return tg.Message{}, nil
}

type mockStorage struct{}

func (m *mockStorage) Get() ([]int64, error) {
	return []int64{1, 2, 3}, nil
}

func (m *mockStorage) Save([]int64) error {
	return nil
}
