package controllers_test

import (
	"log"
	"net"
	"os"
	"testing"

	"github.com/donohutcheon/gowebserver/datalayer/mockdatalayer"
	"github.com/donohutcheon/gowebserver/routes"
	"github.com/donohutcheon/gowebserver/server"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/require"
)

func setup(t *testing.T) (string, *mockdatalayer.MockDataLayer) {
	logger := log.New(os.Stdout, "microservice", log.LstdFlags|log.Lshortfile)

	mockDataLayer := mockdatalayer.New()
	err := mockDataLayer.LoadUserTestData("testdata/users.json")
	require.NoError(t, err)
	err = mockDataLayer.LoadContactTestData("testdata/contacts.json")
	require.NoError(t, err)

	h := routes.NewHandlers(logger, mockDataLayer)
	router := mux.NewRouter()
	h.SetupRoutes(router)

	srv := server.New(router, "", "0")
	l, err := net.Listen("tcp", ":0")
	require.NoError(t, err)

	go func() {
		err := srv.Serve(l)
		require.NoError(t, err)
	} ()

	url := "http://" + l.Addr().String()

	return url, mockDataLayer
}