package main

import (
	"fmt"
	"gitlab.com/donohutcheon/gowebserver/routes"
	"gitlab.com/donohutcheon/gowebserver/server"
	"log"
	"os"
	//"github.com/donohutcheon/gojwt/app"

	"github.com/gorilla/mux"
	_ "github.com/heroku/x/hmetrics/onload"

)

var (
	//CertFile environment variable for CertFile
	CertFile = os.Getenv("CERT_FILE")
	//KeyFile environment variable for KeyFile
	KeyFile = os.Getenv("KEY_FILE")
	//ServiceAddress address to listen on
	ServiceAddress = os.Getenv("SERVICE_ADDRESS")
)

func main() {
	logger := log.New(os.Stdout, "microservice", log.LstdFlags|log.Lshortfile)
	h := routes.NewHandlers(logger)

	router := mux.NewRouter()
	h.SetupRoutes(router)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000" //localhost
	}

	fmt.Println(port)
	srv := server.New(router, ServiceAddress)
	// TODO: Put back in for TLS
	/*err := srv.ListenAndServeTLS(CertFile, KeyFile)*/
	err := srv.ListenAndServe() //Launch the app, visit localhost:8000/api
	if err != nil {
		fmt.Print(err)
	}
}
