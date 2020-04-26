package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/models"
)

func CreateUser(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	if r.Method == http.MethodOptions {
		return nil
	}

	user := models.NewUser(dataLayer)
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		err = errors.Wrap("Invalid request", http.StatusBadRequest, err)
		errors.WriteError(w, err)
		return err
	}

	data, err := user.Create() //Create user
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "User has been created")
	resp.Set("user", data)
	resp.Respond(w)

	return nil
}

// TODO: Move into usersController
func GetCurrentUser(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set( "Access-Control-Allow-Methods", "OPTIONS,GET,POST,PUT,DELETE")
	w.Header().Set("Access-Control-Allow-Headers", "authorization")

	if r.Method == http.MethodOptions {
		return nil
	}
	id := r.Context().Value("userID").(int64)

	user := models.NewUser(dataLayer)
	err := user.GetUser(id)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	user.Password = ""

	resp := response.New(true, "success")
	resp.Set("user", user)

	err = resp.Respond(w)
	if err != nil {
		return err
	}

	return nil
}