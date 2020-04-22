package controllers_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/donohutcheon/gowebserver/datalayer/mockdatalayer"
	"github.com/donohutcheon/gowebserver/routes"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T, route string, handlerFunc routes.HandlerFunc) (string, func()) {
	mux := http.NewServeMux()
	server := httptest.NewServer(mux)
	logger := log.New(os.Stdout, "microservice", log.LstdFlags|log.Lshortfile)
	mockDataLayer := mockdatalayer.New()
	err := mockDataLayer.LoadAccountTestData("testdata/accounts.json")
	require.NoError(t, err)
	err = mockDataLayer.LoadContactTestData("testdata/contacts.json")
	require.NoError(t, err)
	h := routes.NewHandlers(logger, mockDataLayer)
	mux.HandleFunc(route, h.Logger(handlerFunc))
	return server.URL, func() {
		defer server.Close()
	}
}