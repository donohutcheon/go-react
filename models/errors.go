package models

import (
	"net/http"

	e "github.com/donohutcheon/gowebserver/controllers/errors"
)

var (
	ErrLoginFailed = e.NewError("Invalid login credentials", http.StatusForbidden)
	ErrValidationEmail = e.NewError("Email address is required", http.StatusBadRequest)
	ErrValidationPassword = e.NewError("Password is required", http.StatusBadRequest)
	ErrUserDoesNotExist = e.NewError("User does not exist", http.StatusForbidden)
	ErrEmailExists = e.NewError("Email address already in use by another user", http.StatusBadRequest)

	ErrValidationFailed = e.NewError("Invalid request, validation failed", http.StatusBadRequest)

	ErrValidationName  = e.NewError("Contact name is required", http.StatusBadRequest)
	ErrValidationPhone  = e.NewError("Contact phone number is required", http.StatusBadRequest)
)
