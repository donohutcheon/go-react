package controllers

import (
	"encoding/json"
	"gitlab.com/donohutcheon/gowebserver/controllers/errors"
	"log"

	"net/http"

	"gitlab.com/donohutcheon/gowebserver/controllers/response"

	"gitlab.com/donohutcheon/gowebserver/models"
)

func CreateContact(w http.ResponseWriter, r *http.Request, logger *log.Logger) error {
	user := r.Context().Value("user").(int64) //Grab the id of the user that send the request
	contact := &models.Contact{}

	err := json.NewDecoder(r.Body).Decode(contact)
	if err != nil {
		resp := response.New(false, "Error while decoding request body")
		resp.Respond(w)
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	contact.UserID = user
	resp, err := contact.Create()
	if err != nil {
		errors.WriteError(w, err)
		return err
	}
	resp.Respond(w)

	return nil
}

func GetContactsFor (w http.ResponseWriter, r *http.Request, logger *log.Logger) error {
	id := r.Context().Value("user").(int64)
	data := models.GetContacts(id)
	resp := response.New(true, "success")

	err := resp.Set("data", data)
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp.Respond(w)

	return nil
}
