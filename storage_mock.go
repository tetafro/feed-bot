package main

import "time"

type mockStorage struct{}

func (m *mockStorage) GetChats() ([]int64, error) {
	return []int64{1, 2, 3}, nil
}

func (m *mockStorage) SaveChats([]int64) error {
	return nil
}

func (m *mockStorage) GetLastUpdate(feed string) (time.Time, error) {
	return time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), nil
}

func (m *mockStorage) SaveLastUpdate(feed string, t time.Time) error {
	return nil
}
