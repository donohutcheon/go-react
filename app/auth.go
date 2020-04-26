package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"github.com/donohutcheon/gowebserver/datalayer"
)

const AccessTokenLifeSpan = 600
const RefreshTokenLifeSpan = 6000

type JSONWebToken struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

type RefreshJWTReq struct {
	GrantType    string `json:"grantType" sql:"-"`
	RefreshToken string `json:"refreshToken" sql:"-"`
}

type TokenResponse struct {
	ExpiresIn int64 `json:"expiresIn"`
	AccessToken string  `json:"accessToken" sql:"-"`
	RefreshToken string  `json:"refreshToken" sql:"-"`
}

func JwtAuthentication (next http.Handler, logger *log.Logger, dataLayer datalayer.DataLayer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("X-FRAME-OPTIONS", "SAMEORIGIN")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		if r.Header.Get("Access-Control-Request-Headers") != "" {
			w.Header().Set("Access-Control-Allow-Headers", "content-type")
		}
		//w.Header().Set("X-Content-Type-Options","nosniff")


		if r.Method == http.MethodOptions {
			next.ServeHTTP(w, r)
			return
		}
		// List of endpoints that doesn't require auth
		notAuth := []string{
			"/auth/login",
			"/auth/refresh",
			"/auth/sign-up",
			"/status",
		}

		requestPath := r.URL.Path // Current request path

		//check if request does not need authentication, serve the request if it doesn't need it
		for _, value := range notAuth {
			if value == requestPath {
				next.ServeHTTP(w, r)
				return
			}
		}

		//response := response.NewEmptyResponse()
		tokenHeader := r.Header.Get("Authorization") //Grab the token from the header

		if tokenHeader == "" { //Token is missing, returns with error code 403 Unauthorized
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
		tk := &JSONWebToken{}

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

		// TODO: Example, remove
		/*user, err := dataLayer.GetUserByID(tk.UserID)
		if err != nil {
			log.Println(err.Error())
		}
		log.Print(user)*/

		//Everything went well, proceed with the request and set the caller to the user retrieved from the parsed token
		fmt.Printf("User %d", tk.UserID) //Useful for monitoring
		ctx := context.WithValue(r.Context(), "userID", tk.UserID)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r) //proceed in the middleware chain!
	})
}

func CreateToken(userID int64) (*TokenResponse, error){
	token := new(TokenResponse)
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + AccessTokenLifeSpan
	token.ExpiresIn = expireDateTime
	accessToken := &JSONWebToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}

	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	accessTokenString, _ := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	token.AccessToken = accessTokenString

	refreshToken := &JSONWebToken{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: epochSecs + RefreshTokenLifeSpan,
			IssuedAt:  epochSecs,
		},
	}
	signedRefreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	refreshTokenString, _ := signedRefreshToken.SignedString([]byte(os.Getenv("token_password")))
	token.RefreshToken = refreshTokenString

	return token, nil
}

func RefreshToken(rawToken string) (*TokenResponse, error) {
	tk := new(JSONWebToken)

	token, err := jwt.ParseWithClaims(rawToken, tk, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("token_password")), nil
	})
	if err != nil { //Malformed token, returns with http code 403 as usual
		return nil, e.Wrap("Token rejected", http.StatusForbidden, err)
	}

	if !token.Valid { //Token is invalid, maybe not signed on this server
		return nil, e.NewError ("token is not valid", http.StatusForbidden)
	}

	fmt.Printf("UserID %d", tk.UserID)

	//Create JWT token
	tokenResp, err := CreateToken(tk.UserID)
	if err != nil {
		return nil, e.Wrap("token creation failed", http.StatusInternalServerError, err)
	}

	return tokenResp, nil
}
