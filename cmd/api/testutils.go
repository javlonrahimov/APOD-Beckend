package main

import (
	"javlonrahimov/apod/internal/data"
	"javlonrahimov/apod/internal/jsonlog"
	"javlonrahimov/apod/internal/mailer"
	"os"
	"sync"
	"testing"
)

func newTestApplication(t *testing.T) *application {
	return &application{
		config: config{},
		logger: jsonlog.New(os.Stdout, jsonlog.LevelOff, true),
		models: data.Models{},
		mailer: mailer.Mailer{},
		wg:     sync.WaitGroup{},
	}
}