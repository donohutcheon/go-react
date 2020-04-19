package controllers_test

import (
	"github.com/donohutcheon/gowebserver/datalayer/mockdatalayer"
	"log"
	"net/http"
	"net/http/httptest"
	"os"

	"github.com/donohutcheon/gowebserver/routes"
)

func setup(route string, handlerFunc routes.HandlerFunc) (string, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	logger := log.New(os.Stdout, "microservice", log.LstdFlags|log.Lshortfile)
	mockDataLayer := mockdatalayer.New()
	h := routes.NewHandlers(logger, mockDataLayer)
	mux.HandleFunc(route, h.Logger(handlerFunc))
	return server.URL, func() {
		defer server.Close()
	}
}