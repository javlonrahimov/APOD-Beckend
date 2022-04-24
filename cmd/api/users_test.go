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
	handlers := NewHandler(app)
	tests := []struct {
		name         string
		userName     string
		userEmail    string
		userPassword string
		wantStatus   int
		wantCode     int
	}{
		{
			name:         "Valid registration",
			userName:     "Javlon",
			userEmail:    "javlon@gmail.com",
			userPassword: "this_is_a_password",
			wantStatus:   http.StatusAccepted,
			wantCode:     ErrCodeOk,
		},
		{
			name:         "Invalid email",
			userName:     "Javlon",
			userEmail:    "javlongmail",
			userPassword: "this_is_a_password",
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     ErrCodeValidation},
		{
			name:         "Invalid name",
			userName:     "",
			userEmail:    "javlon@gmail.com",
			userPassword: "this_is_a_password",
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     ErrCodeValidation,
		},
		{
			name:         "Invalid password",
			userName:     "Javlon",
			userEmail:    "javlon@gmail.com",
			userPassword: "1234",
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     ErrCodeValidation},
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
			handler := http.HandlerFunc(handlers.Users.Register)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf("handler returned wrong status code: got %v want %v",
					status, tt.wantStatus)
			}

			var m map[string]int
			json.Unmarshal(rr.Body.Bytes(), &m)

			if m["code"] != tt.wantCode {
				t.Errorf("handler returned wrogn code: got %v want %v",
					m["code"], tt.wantCode)
			}
		})
	}
}

func TestUserLogin(t *testing.T) {

	app := newTestApplication()
	handlers := NewHandler(app)

	tests := []struct {
		name         string
		userEmail    string
		userPassword string
		wantStatus   int
		wantCode     int
	}{
		{
			name:         "Valid Login",
			userEmail:    "user1@gmail.com",
			userPassword: "password",
			wantStatus:   http.StatusAccepted,
			wantCode:     ErrCodeOk,
		},
		{
			name:         "Not found email",
			userEmail:    "user@gmail.com",
			userPassword: "password",
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     ErrCodeValidation,
		},
		{
			name:         "Incorrect Password",
			userEmail:    "user1@gmail.com",
			userPassword: "pasword",
			wantStatus:   http.StatusUnprocessableEntity,
			wantCode:     ErrCodeValidation,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"email":"%s","password":"%s"}`, tt.userEmail, tt.userPassword))

			req, err := http.NewRequest("POST", "http:/localhost:4000/v1/users/login", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.Users.Login)
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

func TestUserVerify(t *testing.T) {
	app := newTestApplication()
	handlers := NewHandler(app)

	tests := []struct {
		name       string
		userEmail  string
		userOtp    string
		wantStatus int
		wantCode   int
	}{
		{
			name:       "Valid verify",
			userEmail:  "user1@gmail.com",
			userOtp:    "123456",
			wantStatus: http.StatusAccepted,
			wantCode:   0,
		},
		{
			name:       "Invalid email",
			userEmail:  "user@gmail.com",
			userOtp:    "123456",
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   ErrCodeValidation,
		},
		{
			name:       "Invalid otp",
			userEmail:  "user1@gmail.com",
			userOtp:    "123451",
			wantStatus: http.StatusUnprocessableEntity,
			wantCode:   ErrCodeValidation,
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			var jsonStr = []byte(fmt.Sprintf(`{"email":"%s","otp":"%s"}`, tt.userEmail, tt.userOtp))

			req, err := http.NewRequest("POST", "http:/localhost:4000/v1/users/verify", bytes.NewBuffer(jsonStr))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(handlers.Users.Verify)
			handler.ServeHTTP(rr, req)

			if status := rr.Code; status != tt.wantStatus {
				t.Errorf(rr.Body.String())
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
