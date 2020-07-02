package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/models"
	"github.com/donohutcheon/gowebserver/state"
)

func CreateContact(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	user := r.Context().Value("user").(int64) //Grab the id of the user that send the request
	contact := models.NewContact(state)

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		resp := response.New(false, "Error while decoding request body")
		resp.Respond(w)
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	contact.UserID = user
	data, err := contact.Create()
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "success")
	resp.Set("data", data)
	resp.Respond(w)

	return nil
}

func GetContactsFor(w http.ResponseWriter, r *http.Request, state *state.ServerState) error {
	contact := models.NewContact(state)
	id := r.Context().Value("userID").(int64)
	data, err := contact.GetContacts(id)
	if err == datalayer.ErrNoData {
		return errors.Wrap("User not found", http.StatusNotFound, err)
	} else if err != nil {
		return errors.Wrap("Could not get user's contacts", http.StatusInternalServerError, err)
	}

	resp := response.New(true, "success")
	resp.Set("data", data)

	resp.Respond(w)

	return nil
}
