package feed

import "time"

type mockStorage struct{}

func (m *mockStorage) GetLastUpdate(feed string) (time.Time, error) {
	return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil
}

func (m *mockStorage) SaveLastUpdate(feed string, t time.Time) error {
	return nil
}
