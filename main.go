package main

import (
	"log"
	"os"

	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/donohutcheon/gowebserver/routes"
	"github.com/donohutcheon/gowebserver/server"
	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"
)

var (
	//CertFile environment variable for CertFile
	CertFile = os.Getenv("CERT_FILE")
	//KeyFile environment variable for KeyFile
	KeyFile = os.Getenv("KEY_FILE")
	//ServiceAddress address to listen on
	BindAddress = os.Getenv("BIND_ADDRESS")
	Port = os.Getenv("PORT")
)

func main() {
	logger := log.New(os.Stdout, "microservice", log.LstdFlags|log.Lshortfile)
	dataLayer, err := datalayer.New()
	h := routes.NewHandlers(logger, dataLayer)

	router := mux.NewRouter()
	h.SetupRoutes(router)

	logger.Printf("Server Binding to %s:%s", BindAddress, Port)
	srv := server.New(router, BindAddress, Port)
	// TODO: Put back in for TLS
	/*err := srv.ListenAndServeTLS(CertFile, KeyFile)*/
	err = srv.ListenAndServe() //Launch the app, visit localhost:8000/api
	if err != nil {
		logger.Fatalf("Server failed to start %s", err.Error())
	}
}
