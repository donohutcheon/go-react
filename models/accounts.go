package models

import (
	"database/sql"
	"fmt"
	"github.com/donohutcheon/gowebserver/app"
	"github.com/donohutcheon/gowebserver/datalayer"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/controllers/response"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrLoginFailed = e.NewError("Invalid login credentials", http.StatusForbidden)
	ErrValidationEmail = e.NewError("Email address is required", http.StatusBadRequest)
	ErrValidationPassword = e.NewError("Email address is required", http.StatusBadRequest)
	ErrUserDoesNotExist = e.NewError("User does not exist", http.StatusForbidden)
	ErrEmailExists = e.NewError("Email address already in use by another user", http.StatusBadRequest)
)

type Settings struct {
	ID int `json:"id"`
	ThemeName string `json:"themeName"`
}

//a struct to rep user account
type Account struct {
	datalayer.Model
	Username     string   `json:"username"`
	Email        string   `json:"email"`
	FirstName    string   `json:"firstName"`
	Surname      string   `json:"surname"`
	Age          int      `json:"age"`
	Address      string   `json:"address"`
	Roles        []string `json:"roles"`
	Settings     Settings `json:"settings"`
	Password     string   `json:"password"`
	AccessToken  string   `json:"access_token" sql:"-"`
	RefreshToken string   `json:"refresh_token" sql:"-"`
	dataLayer    datalayer.DataLayer
}

func NewAccount(dataLayer datalayer.DataLayer) *Account {
	account := new(Account)
	account.dataLayer = dataLayer
	return account
}

func (a *Account) convert(account datalayer.Account) {
	a.ID = account.ID
	a.CreatedAt = account.CreatedAt
	a.UpdatedAt = account.UpdatedAt
	a.DeletedAt = account.DeletedAt
	if account.Email.Valid {
		a.Email = account.Email.String
	}
	if account.Password.Valid {
		a.Password = account.Password.String
	}
	a.Roles = []string{"ADMIN","USER"}
	a.Settings.ID = 0
	a.Settings.ThemeName = "default"
}


//Validate incoming user details...
func (a *Account) validate() (response.Response, error) {
	if !strings.Contains(a.Email, "@") {
		return nil, ErrValidationEmail
	}

	if len(a.Password) < 6 {
		return nil, ErrValidationPassword
	}

	//Email must be unique
	//check for errors and duplicate emails
	dl := a.dataLayer
	_, err := dl.GetAccountByEmail(a.Email)
	if err != datalayer.ErrNoData {
		return nil, ErrEmailExists
	}

	return response.New(false, "Requirement passed"), nil
}

func (a *Account) Create() (response.Response, error) {
	resp, err := a.validate()
	if err != nil {
		return resp, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(a.Password), bcrypt.DefaultCost)
	a.Password = string(hashedPassword)

	dl := a.dataLayer
	id, err :=  dl.CreateAccount(a.Email, a.Password)
	if err != nil {
		log.Fatal(err) // TODO: remove
		return nil, err
	}

	dbAccount, err := dl.GetAccountByID(id)
	if err != datalayer.ErrNoData {
		return nil, err
	}

	account := new(Account)
	account.convert(*dbAccount)

	//Create new JWT token for the newly registered account
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + app.AccessTokenLifeSpan
	accessToken := &app.JSONWebToken{
		UserID: dbAccount.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}
	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	accessTokenString, _ := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	account.AccessToken = accessTokenString

	refreshToken := &app.JSONWebToken{
		UserID: account.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: epochSecs + app.RefreshTokenLifeSpan,
			IssuedAt:  epochSecs,
		},
	}
	signedRefreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	refreshTokenString, _ := signedRefreshToken.SignedString([]byte(os.Getenv("token_password")))
	account.RefreshToken = refreshTokenString

	account.Password = "" //delete password

	response := response.New(true, "Account has been created")
	err = response.Set("account", account)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (a *Account) Login(email, password string) (response.Response, error) {
	dataLayer := a.dataLayer
	dbAcc, err := dataLayer.GetAccountByEmail(email)
	if err == sql.ErrNoRows {
		return nil, ErrLoginFailed
	} else if err != nil {
		return nil, err
	}

	a.convert(*dbAcc)

	err = bcrypt.CompareHashAndPassword([]byte(a.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, ErrLoginFailed
	}
	// Worked! Logged In
	a.Password = ""

	// Create JWT token
	var tokenResp app.TokenResponse
	tokenResp, err = app.CreateToken(a.ID)
	if err != nil {
		return nil, e.Wrap("token creation failed", http.StatusInternalServerError, err)
	}

	resp := response.New(true, "Logged In")
	resp["token"] = tokenResp

	return resp, nil
}

func (a *Account) GetAccount(id int64) (error) {
	dl := a.dataLayer
	dbAccount, err := dl.GetAccountByID(id)
	if err == datalayer.ErrNoData {
		return ErrUserDoesNotExist
	} else if err != nil {
		return e.Wrap(fmt.Sprintf("Failed to query account [%d] from database", id), http.StatusInternalServerError, err)
	}

	a.convert(*dbAccount)

	return nil
}
