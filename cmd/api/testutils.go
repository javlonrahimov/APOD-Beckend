package main

import (
	"javlonrahimov/apod/internal/jsonlog"
	"javlonrahimov/apod/internal/mailer"
	"javlonrahimov/apod/internal/mock"
	"net/http/httptest"
	"os"
	"sync"
)

type testServer struct {
	*httptest.Server
}

func newTestApplication() *application {
	return &application{
		config: config{},
		logger: jsonlog.New(os.Stdout, jsonlog.LevelOff, true),
		models: mock.NewModelsMock(),
		mailer: mailer.New("", "", "", "", "", true),
		wg:     sync.WaitGroup{},
	}
}

