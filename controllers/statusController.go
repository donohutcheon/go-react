package controllers

import (
	"github.com/donohutcheon/gowebserver/datalayer"
	"log"
	"net/http"

	"github.com/donohutcheon/gowebserver/controllers/response"
)

func Status(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
	w.Header().Set("X-XSS-Protection", "1; mode=block")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	//w.Header().Set("Sec-Fetch-Site", "same-site")

	resp := response.New(true, "Service is up")
	resp.Respond(w)
	return nil
}
