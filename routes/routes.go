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
type MiddlewareFunc func(next http.Handler, logger *log.Logger, dataLayer datalayer.DataLayer) http.Handler

type Handlers struct {
	logger *log.Logger
	dataLayer datalayer.DataLayer
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (h *Handlers) WrapHandlerFunc(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %v, %s\n", getFunctionName(next), time.Now().Sub(startTime))
		err := next(w, r, h.logger, h.dataLayer)
		if err != nil {
			h.logger.Printf("Controller error: %v", err)
		}
	}
}

func (h *Handlers) WrapMiddlewareFunc(next MiddlewareFunc) mux.MiddlewareFunc {
	return func(mwf http.Handler) http.Handler {
		startTime := time.Now()
		defer h.logger.Printf("request processed in %v, %s\n", getFunctionName(next), time.Now().Sub(startTime))

		return next(mwf, h.logger, h.dataLayer)
	}
}

//SetupRoutes add home route to mux
func (h *Handlers) SetupRoutes(router *mux.Router) {
	router.HandleFunc("/status",  h.WrapHandlerFunc(controllers.Status)).Methods(http.MethodGet)
	router.HandleFunc("/auth/sign-up",  h.WrapHandlerFunc(controllers.CreateUser)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/users/current",  h.WrapHandlerFunc(controllers.GetCurrentUser)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/auth/login", h.WrapHandlerFunc(controllers.Authenticate)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/auth/refresh", h.WrapHandlerFunc(controllers.RefreshToken)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/contacts/new", h.WrapHandlerFunc(controllers.CreateContact)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/me/contacts", h.WrapHandlerFunc(controllers.GetContactsFor)).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/card-transactions/new", h.WrapHandlerFunc(controllers.CreateCardTransaction)).Methods(http.MethodPost, http.MethodOptions)
	router.HandleFunc("/me/card-transactions", h.WrapHandlerFunc(controllers.GetCardTransactions)).Methods(http.MethodGet, http.MethodOptions)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(h.WrapMiddlewareFunc(app.JwtAuthentication)) //attach JWT auth middleware
}

//NewHandlers void
func NewHandlers(logger *log.Logger, dataLayer datalayer.DataLayer) *Handlers {
	return &Handlers{
		logger: logger,
		dataLayer: dataLayer,
	}
}