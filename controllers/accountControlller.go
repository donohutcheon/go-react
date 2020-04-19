package controllers

import (
	"encoding/json"
	"github.com/donohutcheon/gowebserver/datalayer"
	"log"
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/models"
)

func CreateAccount (w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	if r.Method == http.MethodOptions {
		return nil
	}

	account := new(models.Account)
	err := json.NewDecoder(r.Body).Decode(account) //decode the request body into struct and failed if any error occur
	if err != nil {
		err = errors.Wrap("Invalid request", http.StatusBadRequest, err)
		errors.WriteError(w, err)
		return err
	}

	resp, err := account.Create() //Create account
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp.Respond(w)

	return nil
}

// TODO: Move into usersController
func GetCurrentAccount(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set( "Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")

	if r.Method == http.MethodOptions {
		return nil
	}
	id := r.Context().Value("user").(int64)

	account := models.NewAccount(dataLayer)

	err := account.GetAccount(id)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	data, err := dataLayer.GetAccountByID(id)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "success")
	resp.Set("data", data)

	err = resp.Respond(w)
	if err != nil {
		return err
	}

	return nil
}