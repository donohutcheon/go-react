package controllers

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/datalayer"
	"log"
	"net/http"

	"github.com/donohutcheon/gowebserver/app"
	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/models"
)

func Authenticate(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	//w.Header().Set("Sec-Fetch-Site", "same-site")

	if r.Method == http.MethodOptions {
		return nil
	}

	account := models.NewAccount(dataLayer)
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		err = errors.Wrap("Invalid request format", http.StatusBadRequest, err)
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	resp, err := account.Login(account.Email, account.Password)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp.Respond(w)
	return nil
}

func RefreshToken(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	//w.Header().Set("Sec-Fetch-Site", "same-site")

	if r.Method == http.MethodOptions {
		return nil
	}

	refreshTokenReq := new(app.RefreshJWTReq)
	err := json.NewDecoder(r.Body).Decode(refreshTokenReq) //decode the request body into struct and failed if any error occur
	if err != nil {
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	if refreshTokenReq.GrantType != "refresh_token" {
		errors.WriteError(w, errors.NewError("grant type not refresh_token", http.StatusBadRequest))
		return err
	}

	resp, err := app.RefreshToken(refreshTokenReq.RefreshToken)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp.Respond(w)
	return nil
}