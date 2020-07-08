package routes

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/models/auth"
	"github.com/donohutcheon/gowebserver/routes/types"
	"github.com/donohutcheon/gowebserver/state"
	"github.com/gorilla/mux"
)

type Handlers struct {
	serverState *state.ServerState
	staticPath string
	indexPath  string
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func (h *Handlers) WrapHandlerFunc(next types.HandlerFunc) http.HandlerFunc {
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

func (h *Handlers) WrapMiddlewareFunc(next types.MiddlewareFunc, registry map[string]types.RouteEntry) mux.MiddlewareFunc {
	return func(mwf http.Handler) http.Handler {
		startTime := time.Now()
		logger := h.serverState.Logger
		//TODO: Format time
		defer logger.Printf("request processed in %v, %v\n", getFunctionName(next), time.Now().Sub(startTime))

		return next(mwf, h.serverState, registry)
	}
}

//SetupRoutes add home route to mux
func (h *Handlers) SetupRoutes(router *mux.Router) {
	registry := types.GetRouteRegistry()
	for r, e := range registry {
		if e.Handler == nil {
			continue
		}
		router.HandleFunc(r, h.WrapHandlerFunc(e.Handler)).Methods(e.Methods...)
	}

	router.Use(mux.CORSMethodMiddleware(router))
	router.Use(CORSAccessControlAllowOrigin(router))
	router.Use(h.WrapMiddlewareFunc(JwtAuthentication, registry)) //attach JWT auth middleware
}

func CORSAccessControlAllowOrigin(r *mux.Router) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if req.Method == http.MethodOptions {
				w.Header().Set("Access-Control-Allow-Origin", "*")
			}

			next.ServeHTTP(w, req)
		})
	}
}

//NewHandlers void
func NewHandlers(state *state.ServerState) *Handlers {
	return &Handlers{
		serverState: state,
		staticPath: "static/build/",
		indexPath: "index.html",
	}
}

func JwtAuthentication (next http.Handler, state *state.ServerState, registry map[string]types.RouteEntry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if r.Header.Get("Access-Control-Request-Headers") != "" {
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		}

		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}

		requestPath := r.URL.Path // Current request path
		state.Logger.Printf(requestPath)
		for k, v := range registry {
			state.Logger.Printf("%s %+v", k, v)
		}

		//check if request does not need authentication, serve the request if it doesn't need it
		var isPublicMatch bool
		err := state.Router.Walk(func (route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
			pathTemplate, err := route.GetPathTemplate()
			if err != nil {
				return err
			}

			pathRegexp, err := route.GetPathRegexp()
			if err != nil {
				return err
			}
			matched, err := regexp.MatchString(pathRegexp, requestPath)
			if err != nil {
				return err
			}

			isAPI, err := regexp.MatchString("^/api/", requestPath)
			if err != nil {
				return err
			}

			if v, ok := registry[pathTemplate]; matched && ok {
				if pathTemplate == "/" && isAPI {
					return nil
				}
				isPublicMatch = v.Public
				return nil
			}

			return nil
		})
		if err != nil {
			resp := response.New(false, "Internal server error.  Routing failed")
			w.WriteHeader(http.StatusInternalServerError)
			w.Header().Add("Content-Type", "application/json")
			resp.Respond(w)
			return
		}
		if isPublicMatch {
			next.ServeHTTP(w, r)
			return
		}

		tokenHeader := r.Header.Get("Authorization")
		if tokenHeader == "" {
			resp := response.New(false, "Missing auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Respond(w)
			return
		}

		splitted := strings.Split(tokenHeader, " ") //The token normally comes in format `Bearer {token-body}`, we check if the retrieved token matched this requirement
		if len(splitted) != 2 {
			resp := response.New(false, "Invalid/Malformed auth token")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Respond(w)
			return
		}

		tokenPart := splitted[1] //Grab the token part, what we are truly interested in
		tk := &auth.JSONWebToken{}

		token, err := jwt.ParseWithClaims(tokenPart, tk, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("token_password")), nil
		})

		if err != nil { //Malformed token, returns with http code 403 as usual
			message := fmt.Sprintf("Token rejected, %s", err.Error())
			resp := response.New(false, message)
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Respond(w)
			return
		}

		if !token.Valid { //Token is invalid, maybe not signed on this server
			resp := response.New(false, "Token is not valid.")
			w.WriteHeader(http.StatusForbidden)
			w.Header().Add("Content-Type", "application/json")
			resp.Respond(w)
			return
		}

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %d", tk.UserID) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "userID", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}


