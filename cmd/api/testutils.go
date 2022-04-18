package main

import (
	"javlonrahimov/apod/internal/jsonlog"
	"javlonrahimov/apod/internal/mailer"
	"javlonrahimov/apod/internal/mock"
	"net/http"
	"net/http/httptest"
	"os"
	"sync"
	"testing"
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

func newTestServer(t *testing.T, h http.Handler) *testServer {
	ts := httptest.NewTLSServer(h)
	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}
	return &testServer{ts}
}

func executeRequest(req *http.Request, app *application) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.routes().ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
