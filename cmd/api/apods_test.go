package main

import (
	"encoding/json"
	"fmt"
	"javlonrahimov/apod/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApodGetAll(t *testing.T) {
	app := newTestApplication()
	handlers := NewHandler(app)
	tests := []struct {
		name       string
		q          string
		filtres    data.Filters
		wantStatus int
		wantCode   int
	}{
		{
			name: "Get all ok",
			q:    "",
			filtres: data.Filters{
				Page:     1,
				PageSize: 10,
			},
			wantStatus: http.StatusOK,
			wantCode:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			req, err := http.NewRequest("GET", fmt.Sprintf("http:/localhost:4000/v1/apods?q=%s&page=%d&page_size=%d&sort=%s", tt.q, tt.filtres.Page, tt.filtres.PageSize, tt.filtres.Sort), nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.Apods.GetAll)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}

			var m map[string]int
			json.Unmarshal(rr.Body.Bytes(), &m)

			if m["code"] != tt.wantCode {
				t.Errorf("handler returned unexpected body: got %v want %v",
					m["code"], tt.wantCode)
			}

		})
	}
}
