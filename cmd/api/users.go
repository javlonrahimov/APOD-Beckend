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

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &data.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}
	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateUser(v, user); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrDuplicateEmail):
			v.AddError("email", "a user with this already exists")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, data.ApodsRead)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	otp, err := app.models.Otps.New(user.ID, 10*time.Minute)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	app.background(func() {
		// data := []byte(otp.Plaintext)

		_, err := http.Get(fmt.Sprintf(
			"https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%v",
			app.config.telegram.botToken, app.config.telegram.channelId, otp.Plaintext,
		))

		// err = app.mailer.Send([]string{user.Email}, data)
		if err != nil {
			app.logger.PrintError(err, nil)
		}
	})

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) loginUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		UserEmail string `json:"email"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateEmail(v, input.UserEmail); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidatePasswordPlaintext(v, input.Password); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.UserEmail)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no user found with this email")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if !user.Activated {
		app.unactivatedAccountResponse(w, r)
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !match {
		app.invalidCredentialsResponse(w, r)
		return
	}

	accessToken, err := app.models.Tokens.New(user.ID, data.AccessTokenExpire, data.ScopeAccess)
	if err != nil {
		return
	}
	refreshToken, err := app.models.Tokens.New(user.ID, data.RefreshTokenExpire, data.ScopeRefresh)

	tokenPair := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken.Plaintext,
		RefreshToken: refreshToken.Plaintext,
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user, "tokens": tokenPair}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) verifyUserHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		UserEmail    string `json:"email"`
		OtpPlainText string `json:"otp"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	v := validator.New()

	if data.ValidateOtpPlaintext(v, input.OtpPlainText); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	if data.ValidateEmail(v, input.UserEmail); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	user, err := app.models.Users.GetByEmail(input.UserEmail)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("email", "no user found with this email")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	otp, err := app.models.Otps.GetForUser(user.ID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			v.AddError("otp", "incarrect otp")
			app.failedValidationResponse(w, r, v.Errors)
		case errors.Is(err, data.ErrOtpExpired):
			v.AddError("otp", "otp expired")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = bcrypt.CompareHashAndPassword(otp.Hash, []byte(input.OtpPlainText))
	if err != nil {
		switch {
		case errors.Is(err, bcrypt.ErrMismatchedHashAndPassword):
			v.AddError("otp", "incorrect otp")
			app.failedValidationResponse(w, r, v.Errors)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	user.Activated = true
	app.models.Users.Update(user)

	accessToken, err := app.models.Tokens.New(user.ID, data.AccessTokenExpire, data.ScopeAccess)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	refreshToken, err := app.models.Tokens.New(user.ID, data.RefreshTokenExpire, data.ScopeRefresh)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	tokenPair := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
	}{
		AccessToken:  accessToken.Plaintext,
		RefreshToken: refreshToken.Plaintext,
	}

	err = app.writeJSON(w, http.StatusAccepted, envelope{"user": user, "tokens": tokenPair}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Otps.DeleteAllForUser(user.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
