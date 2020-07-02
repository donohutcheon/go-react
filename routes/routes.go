package routes

import (
	"net/http"
	"reflect"
	"runtime"
	"time"

	"github.com/donohutcheon/gowebserver/app"
	"github.com/donohutcheon/gowebserver/controllers"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/gorilla/mux"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request, handlerState *state.ServerState) error
type MiddlewareFunc func(next http.Handler, state *state.ServerState) http.Handler

type Handlers struct {
	serverState *state.ServerState
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (h *Handlers) WrapHandlerFunc(next HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		logger := h.serverState.Logger
		//TODO: Format time
		defer logger.Printf("request processed in %v, %v\n", getFunctionName(next),  time.Now().Sub(startTime))
		err := next(w, r, h.serverState)
		if err != nil {
			logger.Printf("Controller error: %v", err)
		}
	}
}

func (h *Handlers) WrapMiddlewareFunc(next MiddlewareFunc) mux.MiddlewareFunc {
	return func(mwf http.Handler) http.Handler {
		startTime := time.Now()
		logger := h.serverState.Logger
		//TODO: Format time
		defer logger.Printf("request processed in %v, %v\n", getFunctionName(next), time.Now().Sub(startTime))

		return next(mwf, h.serverState)
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
	router.HandleFunc("/users/confirm/{nonce}", h.WrapHandlerFunc(controllers.ConfirmUserSignUp)).Methods(http.MethodGet, http.MethodOptions)
	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(h.WrapMiddlewareFunc(app.JwtAuthentication)) //attach JWT auth middleware
}

//NewHandlers void
func NewHandlers(state *state.ServerState) *Handlers {
	return &Handlers{
		serverState: state,
	}
}