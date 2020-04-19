package routes

import (
	"log"
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/donohutcheon/gowebserver/app"
	"github.com/donohutcheon/gowebserver/controllers"
	"github.com/donohutcheon/gowebserver/datalayer"
	"github.com/gorilla/mux"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, logger *log.Logger, dataLayer datalayer.DataLayer) error

type Handlers struct {
	logger *log.Logger
	dataLayer datalayer.DataLayer
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (h *Handlers) Logger(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %v, %s\n", getFunctionName(next), time.Now().Sub(startTime))
		err := next(w, r, h.logger, h.dataLayer)
		if err != nil {
			h.logger.Printf("Controller error: %v", err)
		}
	}
}

//SetupRoutes add home route to mux
func (h *Handlers) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/status",  h.Logger(controllers.Status)).Methods(http.MethodGet)
	router.HandleFunc("/auth/sign-up",  h.Logger(controllers.CreateAccount)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/users/current",  h.Logger(controllers.GetCurrentAccount)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/login", h.Logger(controllers.Authenticate)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth/refresh", h.Logger(controllers.RefreshToken)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/contacts/new", h.Logger(controllers.CreateContact)).Methods(http.MethodPost)
	router.HandleFunc("/me/contacts", h.Logger(controllers.GetContactsFor)).Methods(http.MethodGet) //  user/2/contacts
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(app.JwtAuthentication) //attach JWT auth middleware
}

//NewHandlers void
func NewHandlers(logger *log.Logger, dataLayer datalayer.DataLayer) *Handlers {
	return &Handlers{
		logger: logger,
		dataLayer: dataLayer,
	}
}