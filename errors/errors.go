package errors

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)


type httpWrappedError interface {
	Error() string
	String() string
	ErrorCode() int
	ErrorName() string
	JSON() []byte
}

// HTTPError is a base type for errors that should return a structured JSON object
type HTTPError struct {
	Code     int                    `json:"code"`
	Name     string                 `json:"name"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}

// BadRequestError is a HTTPError that indicates a 400 status code
type BadRequestError struct {
	*HTTPError
}

// AuthorizationError is a HTTPError that indicates a 401 status code
type AuthorizationError struct {
	*HTTPError
}

// NewBadRequestError creates a new BadRequestError
func NewBadRequestError(message string, metadata map[string]interface{}) *BadRequestError {
	e := NewHTTPError(http.StatusBadRequest, message, metadata)
	return &BadRequestError{e}
}

// NewAuthorizationError creates a new AuthorizationError
func NewAuthorizationError(message string, metadata map[string]interface{}) *AuthorizationError {
	e := NewHTTPError(http.StatusUnauthorized, message, metadata)
	return &AuthorizationError{e}
}

// NewHTTPError creates a new HTTPError
func NewHTTPError(code int, message string, metadata map[string]interface{}) *HTTPError {
	return &HTTPError{
		Code:     code,
		Message:  message,
		Metadata: metadata,
		Name:     http.StatusText(code),
	}
}

// Error outputs the error as stringified JSON
func (e *HTTPError) Error() string {
	return string(e.JSON())
}

// String formats the error as a string
func (e *HTTPError) String() string {
	return fmt.Sprintf("%s (%d): %s, %v", e.Name, e.Code, e.Message, e.Metadata)
}

// JSON outputs the error to JSON
func (e *HTTPError) JSON() []byte {
	val, _ := json.Marshal(e)

	return val
}

// ErrorCode gets the error code
func (e *HTTPError) ErrorCode() int {
	return e.Code
}

// ErrorName gets the error name
func (e *HTTPError) ErrorName() string {
	return e.Name
}

// ErrorWriter does the work of generating a HTTP response for a httpWrappedError
func ErrorWriter(err httpWrappedError, rw http.ResponseWriter) {
	j := err.JSON()

	contentLength := strconv.Itoa(len(j))
	log.WithFields(log.Fields{
		"ContentLength": contentLength,
		"ErrorCode":     err.ErrorCode(),
		"Error":         err,
	}).Infof("Logging %s", err.ErrorName())

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Content-Length", contentLength)
	rw.WriteHeader(err.ErrorCode())
	_, writeErr := rw.Write(j)
	if writeErr != nil {
		log.Error(writeErr)
	}
}
