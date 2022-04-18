package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserRegistration(t *testing.T) {
	app := newTestApplication()
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		wantStatus   int
		wantCode     int
	}{
		{"Valid registration", "Javlon", "javlon@gmail.com", "this_is_a_password", http.StatusAccepted, ErrCodeOk},
		{"Invalid email", "Javlon", "javlongmail", "this_is_a_password", http.StatusUnprocessableEntity, ErrCodeValidation},
		{"Invalid name", "", "javlon@gmail.com", "this_is_a_password",http.StatusUnprocessableEntity, ErrCodeValidation},
		{"Invalid password", "Javlon", "javlon@gmail.com", "1234", http.StatusUnprocessableEntity, ErrCodeValidation},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"name":"%s","email":"%s","password":"%s"}`, tt.userName, tt.userEmail, tt.userPassword))

			req, err := http.NewRequest("POST", "http:/localhost:4000/v1/users", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.registerUserHandler)
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
