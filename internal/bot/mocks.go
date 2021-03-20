package bot

import tg "github.com/go-telegram-bot-api/telegram-bot-api"

type mockStorage struct{}

func (m *mockStorage) GetChats() ([]int64, error) {
	return []int64{1, 2, 3}, nil
}

func (m *mockStorage) SaveChats([]int64) error {
	return nil
}

type mockAPI struct{}

func (m *mockAPI) GetUpdatesChan(tg.UpdateConfig) (tg.UpdatesChannel, error) {
	return nil, nil
}

func (m *mockAPI) Send(tg.Chattable) (tg.Message, error) {
	return tg.Message{}, nil
}
