package main

import (
	"javlonrahimov/apod/internal/data"
	"javlonrahimov/apod/internal/jsonlog"
	"javlonrahimov/apod/internal/mailer"
	"sync"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		config: config{},
		logger: jsonlog.New(nil, jsonlog.LevelOff),
		models: data.Models{},
		mailer: mailer.Mailer{},
		wg:     sync.WaitGroup{},
	}
}