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
	"github.com/donohutcheon/gowebserver/app"
	e "github.com/donohutcheon/gowebserver/controllers/errors"
	"github.com/donohutcheon/gowebserver/datalayer"
	"golang.org/x/crypto/bcrypt"
)

type Settings struct {
	ID int `json:"id"`
	ThemeName string `json:"themeName"`
}

type User struct {
	datalayer.Model
	dataLayer    datalayer.DataLayer
	Email        string    `json:"email"`
	FirstName    string    `json:"firstName"`
	Surname      string    `json:"surname"`
	Age          int       `json:"age"`
	Address      string    `json:"address"`
	Roles        []string  `json:"roles"`
	Settings     Settings  `json:"settings"`
	Password     string    `json:"password,omitempty"`
	AccessToken  string    `json:"accessToken,omitempty" sql:"-"`
	RefreshToken string    `json:"refreshToken,omitempty" sql:"-"`
	LoggedOutAt  time.Time `json:"loggedOutAt,omitempty"`
}

func NewUser(dataLayer datalayer.DataLayer) *User {
	user := new(User)
	user.dataLayer = dataLayer
	return user
}

func (u *User) convert(user datalayer.User) {
	u.ID = user.ID
	u.CreatedAt = user.CreatedAt
	u.UpdatedAt = user.UpdatedAt
	u.DeletedAt = user.DeletedAt
	if user.Email.Valid {
		u.Email = user.Email.String
	}
	if user.Password.Valid {
		u.Password = user.Password.String
	}
	u.Roles = []string{"ADMIN","USER"}
	u.Settings.ID = 0
	u.Settings.ThemeName = "default"
}

//Validate incoming user details...
func (u *User) validate() error {
	if !strings.Contains(u.Email, "@") {
		return ErrValidationEmail
	}

	if len(u.Password) < 6 {
		return ErrValidationPassword
	}

	//Email must be unique
	//check for errors and duplicate emails
	dl := u.dataLayer
	_, err := dl.GetUserByEmail(u.Email)
	if err != sql.ErrNoRows {
		return ErrEmailExists
	}

	return nil
}

func (u *User) Create() (*User, error) {
	err := u.validate()
	if err != nil {
		return nil, err
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	u.Password = string(hashedPassword)

	dl := u.dataLayer
	// TODO: Include roles
	id, err :=  dl.CreateUser(u.Email, u.Password)
	if err != nil {
		log.Fatal(err) // TODO: remove
		return nil, err
	}

	dbUser, err := dl.GetUserByID(id)
	if err != nil && err != datalayer.ErrNoData {
		return nil, err
	}

	user := new(User)
	user.convert(*dbUser)

	//Create new JWT token for the newly registered user
	now := time.Now()
	epochSecs := now.Unix()
	expireDateTime := epochSecs + app.AccessTokenLifeSpan
	accessToken := &app.JSONWebToken{
		UserID: dbUser.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expireDateTime,
			IssuedAt:  epochSecs,
		},
	}
	signedAccessToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), accessToken)
	accessTokenString, _ := signedAccessToken.SignedString([]byte(os.Getenv("token_password")))
	user.AccessToken = accessTokenString

	refreshToken := &app.JSONWebToken{
		UserID: user.ID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: epochSecs + app.RefreshTokenLifeSpan,
			IssuedAt:  epochSecs,
		},
	}
	signedRefreshToken := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), refreshToken)
	refreshTokenString, _ := signedRefreshToken.SignedString([]byte(os.Getenv("token_password")))
	user.RefreshToken = refreshTokenString

	user.Password = "" //delete password

	return user, nil
}

func (u *User) Login(email, password string) (*app.TokenResponse, error) {
	dataLayer := u.dataLayer
	dbAcc, err := dataLayer.GetUserByEmail(email)
	if err == sql.ErrNoRows {
		return nil, ErrLoginFailed
	} else if err != nil {
		return nil, err
	}

	u.convert(*dbAcc)

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return nil, ErrLoginFailed
	}
	// Worked! Logged In
	u.Password = ""

	// Create JWT token
	tokenResp, err := app.CreateToken(u.ID)
	if err != nil {
		return nil, e.Wrap("token creation failed", http.StatusInternalServerError, err)
	}

	return tokenResp, nil
}

func (u *User) GetUser(id int64) (error) {
	dl := u.dataLayer
	dbUser, err := dl.GetUserByID(id)
	if err == datalayer.ErrNoData {
		return ErrUserDoesNotExist
	} else if err != nil {
		return e.Wrap(fmt.Sprintf("Failed to query user [%d] from database", id), http.StatusInternalServerError, err)
	}

	u.convert(*dbUser)

	return nil
}
