package apperr

import "net/http"

type Code string

const (
	CodeValidation          Code = "VALIDATION_ERROR"
	CodeSecretNotFound      Code = "SECRET_NOT_FOUND"
	CodeSecretExpired       Code = "SECRET_EXPIRED"
	CodeSecretAlreadyBurned Code = "SECRET_ALREADY_BURNED"
	CodeDestroyTokenInvalid Code = "DESTROY_TOKEN_INVALID"
	CodeSessionRequired     Code = "SESSION_REQUIRED"
	CodeAdminRequired       Code = "ADMIN_AUTH_REQUIRED"
	CodeInternal            Code = "INTERNAL_ERROR"
)

type Error struct {
	Code    Code
	Message string
	Status  int
	Details any
}

func (e *Error) Error() string { return e.Message }

func New(status int, code Code, message string) *Error {
	return &Error{Status: status, Code: code, Message: message}
}

func Validation(message string) *Error      { return New(http.StatusBadRequest, CodeValidation, message) }
func NotFound(message string) *Error        { return New(http.StatusNotFound, CodeSecretNotFound, message) }
func Gone(code Code, message string) *Error { return New(http.StatusGone, code, message) }
func Unauthorized(code Code, message string) *Error {
	return New(http.StatusUnauthorized, code, message)
}
func InvalidDestroyToken() *Error {
	return New(http.StatusBadRequest, CodeDestroyTokenInvalid, "Destroy token is invalid.")
}
func Internal() *Error {
	return New(http.StatusInternalServerError, CodeInternal, "Internal server error.")
}
