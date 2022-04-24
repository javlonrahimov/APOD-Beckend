package main

import (
	"errors"
	"fmt"
	"javlonrahimov/apod/internal/data"
	"javlonrahimov/apod/internal/validator"
	"net/http"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type UserHandler interface {
	Register(w http.ResponseWriter, r *http.Request)
	Login(w http.ResponseWriter, r *http.Request)
	Verify(w http.ResponseWriter, r *http.Request)
}

type userApi struct {
	app *application
}

func NewUserApi(app *application) UserHandler {
	return &userApi{app: app}
}

func (u *userApi) Register(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := u.app.readJSON(w, r, &input)
	if err != nil {
		u.app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		u.app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = u.app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this already exists")
			u.app.failedValidationResponse(w, r, v.Errors)
		default:
			u.app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = u.app.models.Permissions.AddForUser(user.ID, data.ApodsRead)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}

	otp, err := u.app.models.Otps.New(user.ID, 10*time.Minute)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}

	u.app.background(func() {
		// data := []byte(otp.Plaintext)

		_, err := http.Get(fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%v",
			u.app.config.telegram.botToken, u.app.config.telegram.channelId, otp.Plaintext,
		))

		// err = app.mailer.Send([]string{user.Email}, data)
		if err != nil {
			u.app.logger.PrintError(err, nil)
		}
	})

	err = u.app.writeJSON(w, http.StatusAccepted, 0, envelope{"user": user}, nil)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
	}

}

func (u *userApi) Login(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserEmail string `json:"email"`
		Password  string `json:"password"`
	}

	err := u.app.readJSON(w, r, &input)
	if err != nil {
		u.app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.UserEmail); !v.Valid() {
		u.app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidatePasswordPlaintext(v, input.Password); !v.Valid() {
		u.app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := u.app.models.Users.GetByEmail(input.UserEmail)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no user found with this email")
			u.app.failedValidationResponse(w, r, v.Errors)
		default:
			u.app.serverErrorResponse(w, r, err)
		}
		return
	}

	if !user.Activated {
		u.app.inactiveAccountResponse(w, r)
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		u.app.invalidCredentialsResponse(w, r)
		return
	}

	accessToken, err := u.app.models.Tokens.New(user.ID, data.AccessTokenExpire, data.ScopeAccess)
	if err != nil {
		return
	}
	refreshToken, err := u.app.models.Tokens.New(user.ID, data.RefreshTokenExpire, data.ScopeRefresh)

	tokenPair := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken.Plaintext,
		RefreshToken: refreshToken.Plaintext,
	}

	err = u.app.writeJSON(w, http.StatusAccepted, ErrCodeOk, envelope{"user": user, "tokens": tokenPair}, nil)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
	}
}

func (u *userApi) Verify(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserEmail    string `json:"email"`
		OtpPlainText string `json:"otp"`
	}

	err := u.app.readJSON(w, r, &input)
	if err != nil {
		u.app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateOtpPlaintext(v, input.OtpPlainText); !v.Valid() {
		u.app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidateEmail(v, input.UserEmail); !v.Valid() {
		u.app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := u.app.models.Users.GetByEmail(input.UserEmail)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no user found with this email")
			u.app.failedValidationResponse(w, r, v.Errors)
		default:
			u.app.serverErrorResponse(w, r, err)
		}
		return
	}

	otp, err := u.app.models.Otps.GetForUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("otp", "incorrect otp")
			u.app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrOtpExpired):
			v.AddError("otp", "otp expired")
			u.app.failedValidationResponse(w, r, v.Errors)
		default:
			u.app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(otp.Hash, []byte(input.OtpPlainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			v.AddError("otp", "incorrect otp")
			u.app.failedValidationResponse(w, r, v.Errors)
		default:
			u.app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	u.app.models.Users.Update(user)

	accessToken, err := u.app.models.Tokens.New(user.ID, data.AccessTokenExpire, data.ScopeAccess)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}
	refreshToken, err := u.app.models.Tokens.New(user.ID, data.RefreshTokenExpire, data.ScopeRefresh)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}

	tokenPair := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken.Plaintext,
		RefreshToken: refreshToken.Plaintext,
	}

	err = u.app.writeJSON(w, http.StatusAccepted, 0, envelope{"user": user, "tokens": tokenPair}, nil)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
		return
	}

	err = u.app.models.Otps.DeleteAllForUser(user.ID)
	if err != nil {
		u.app.serverErrorResponse(w, r, err)
	}
}
