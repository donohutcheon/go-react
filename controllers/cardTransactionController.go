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

func CreateCardTransaction(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	if r.Method == http.MethodOptions {
		return nil
	}

	userID := r.Context().Value("userID").(int64) //Grab the id of the userID that send the request
	cardTransaction := models.NewCardTransaction(&dataLayer)

	err := json.NewDecoder(r.Body).Decode(cardTransaction)
	if err != nil {
		resp := response.New(false, "Error while decoding request body")
		resp.Respond(w)
		errors.WriteError(w, err, http.StatusBadRequest)
		return err
	}

	cardTransaction.UserID = userID
	data, err := cardTransaction.CreateCardTransaction()
	if err != nil {
		errors.WriteError(w, err)
		return err
	}

	resp := response.New(true, "success")
	resp.Set("cardTransaction", data)

	resp.Respond(w)

	return nil
}

func GetCardTransactions(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	if r.Method == http.MethodOptions {
		return nil
	}

	cardTransaction := models.NewCardTransaction(&dataLayer)
	userID := r.Context().Value("userID").(int64)
	data, err := cardTransaction.GetCardTransactionsByUserID(userID)
	if err != nil && err != datalayer.ErrNoData {
		errors.WriteError(w, err, http.StatusInternalServerError)
		return err
	}

	resp := response.New(true, "success")
	resp.Set("cardTransactions", data)

	resp.Respond(w)

	return nil
}
