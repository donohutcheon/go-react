package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/donohutcheon/gowebserver/app"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/datalayer"
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

	user := models.NewUser(dataLayer)
	err := json.NewDecoder(r.Body).Decode(user)
	if err != nil {
		err = errors.Wrap("Invalid request format", http.StatusBadRequest, err)
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	data, err := user.Login(user.Email, user.Password)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "Logged In")
	resp["token"] = data
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

	data, err := app.RefreshToken(refreshTokenReq.RefreshToken)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "Tokens refreshed")
	resp.Set("token", data)
	resp.Respond(w)

	return nil
}