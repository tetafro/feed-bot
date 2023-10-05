package main

import (
	"context"
	"log"
)

// LogNotifier prints all incoming messages. This is a notifier for
// debugging purposes.
type LogNotifier struct{}

// NewLogNotifier create new log notifier.
func NewLogNotifier() *LogNotifier {
	return &LogNotifier{}
}

// Notify prints item to stdout.
func (ln *LogNotifier) Notify(_ context.Context, item Item) error {
	log.Printf("[notify] New item: %s", item)
	return nil
}