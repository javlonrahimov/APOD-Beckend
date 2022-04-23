package main

import (
	"fmt"
	"javlonrahimov/apod/internal/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestApodGetAll(t *testing.T) {
	app := newTestApplication()
	tests := []struct{
		name string
		q string
		filtres data.Filters
		wantStatus int
		wantCode int
	} {
		{
			name:       "",
			q:          "",
			filtres:    data.Filters{},
			wantStatus: 0,
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
			handler := http.HandlerFunc(app.loginUserHandler)
			handler.ServeHTTP(rr, req)

			


		})
	}
}