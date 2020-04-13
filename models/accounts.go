package models

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	e "gitlab.com/donohutcheon/gowebserver/controllers/errors"
	"gitlab.com/donohutcheon/gowebserver/controllers/response"
	"golang.org/x/crypto/bcrypt"
)

const AccessTokenLifeSpan = 600
const RefreshTokenLifeSpan = 6000

var ErrLoginFailed = e.NewError("Invalid login credentials. Please try again", http.StatusForbidden)
var ErrValidationEmail = e.NewError("Email address is required", http.StatusBadRequest)
var ErrValidationPassword = e.NewError("Email address is required", http.StatusBadRequest)
var ErrUserDoesNotExist = e.NewError("User does not exist", http.StatusForbidden)
var ErrEmailExists = e.NewError("Email address already in use by another user", http.StatusBadRequest)

/*
JWT claims struct
*/
type JSONWebToken struct {
	UserID int64 `json:"userID"`
	jwt.StandardClaims
}

type Settings struct {
	ID int `json:"id"`
	ThemeName string `json:"themeName"`
}

//a struct to rep user account
type Account struct {
	Model
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	FirstName   string   `json:"firstName"`
	Surname     string   `json:"surname"`
	Age         int      `json:"age"`
	Address     string   `json:"address"`
	Roles       []string `json:"roles"`
	Settings    Settings `json:"settings"`
	Password    string   `json:"password"`
	AccessToken string   `json:"access_token" sql:"-"`
	RefreshToken string  `json:"refresh_token" sql:"-"`
}

type Token struct {
	ExpiresIn int64 `json:"expires_in"`
	AccessToken string  `json:"access_token" sql:"-"`
	RefreshToken string  `json:"refresh_token" sql:"-"`
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
	var id sql.NullString
	GetConn().QueryRow("select id from accounts where email = $1", a.Email).Scan(&id)
	if id.Valid {
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

	//GetDB().Create(account)
	result, err := GetConn().Exec("insert into accounts(email, password) values (?, ?)", a.Email, a.Password)
	if err != nil {
		log.Fatal(err) // TODO: remove
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err) // TODO: remove
		return nil, err
	}

	account, err := GetUser(id)
	if err != nil {
		return nil, err
	}

	//Create new JWT token for the newly registered account
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + AccessTokenLifeSpan
	accessToken := &JSONWebToken{
		account.ID,
		jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}
	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	accessTokenString, _ := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	account.AccessToken = accessTokenString

	refreshToken := &JSONWebToken{
		account.ID,
		jwt.StandardClaims{
			ExpiresAt: epochSecs + RefreshTokenLifeSpan,
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

func Login(email, password string) (response.Response, error) {
	account := &Account{}
	token := &Token{}
	err := GetConn().QueryRow("select id, email, password, token, created_at, updated_at, deleted_at from accounts where email = ?", email).Scan(&account.ID, &account.Email, &account.Password, &account.AccessToken, &account.CreatedAt, &account.UpdatedAt, &account.DeletedAt)
	if err == sql.ErrNoRows {
		return nil, ErrLoginFailed
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, ErrLoginFailed
	}
	//Worked! Logged In
	account.Password = ""

	//Create JWT token
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + AccessTokenLifeSpan
	token.ExpiresIn = expireDateTime
	accessToken := &JSONWebToken{
		account.ID,
		jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}

	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	accessTokenString, _ := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	token.AccessToken = accessTokenString

	refreshToken := &JSONWebToken{
		account.ID,
		jwt.StandardClaims{
			ExpiresAt: epochSecs + RefreshTokenLifeSpan,
			IssuedAt:  epochSecs,
		},
	}
	signedRefreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	refreshTokenString, _ := signedRefreshToken.SignedString([]byte(os.Getenv("token_password")))
	token.RefreshToken = refreshTokenString

	resp := response.New(true, "Logged In")
	resp["token"] = token

	return resp, nil
}

func GetUser(id int64) (*Account, error) {
	account := &Account{}
	var token sql.NullString
	var createdAt, updatedAt, deletedAt sql.NullTime
	err := GetConn().QueryRow(`SELECT id, email, password, token, created_at, updated_at, deleted_at FROM accounts WHERE id=?`, id).Scan(&account.ID, &account.Email, &account.Password, &token, &createdAt, &updatedAt, &deletedAt)
	if err == sql.ErrNoRows {
		fmt.Println(false, "Invalid login credentials. Please try again")
		return nil, ErrUserDoesNotExist
	} else if err != nil {
		fmt.Printf("Failed to query account [%d] from database", id) // TODO: remove
		return nil, e.Wrap(fmt.Sprintf("Failed to query account [%d] from database", id), http.StatusInternalServerError, err)
	}

	if token.Valid {
		account.AccessToken = token.String
	}
	if createdAt.Valid {
		account.CreatedAt = &createdAt.Time
	}
	if updatedAt.Valid {
		account.UpdatedAt = &updatedAt.Time
	}
	if deletedAt.Valid {
		account.DeletedAt = &deletedAt.Time
	}

	account.Username = account.Email
	account.Password = ""
	account.Roles = []string{"ADMIN","USER"}
	account.Settings.ID = 0
	account.Settings.ThemeName = "default"

	return account, nil
}
