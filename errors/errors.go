package errors

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

type HTTPWrappedError interface {
	Error() string
	String() string
	ErrorCode() int
	ErrorName() string
	Json() []byte
}

type HTTPError struct {
	Code     int                    `json:"code"`
	Name     string                 `json:"name"`
	Message  string                 `json:"message"`
	Metadata map[string]interface{} `json:"metadata"`
}

type RequestError struct {
	*HTTPError
}

type AuthorizationError struct {
	*HTTPError
}

type UnsupportedAlgorithmError struct {
	Message   string `json:"message"`
	Algorithm string `json:"algorithm"`
	Name      string `json:"name"`
}

func NewUnsupportedAlgorithmError(message string, algorithm string) *UnsupportedAlgorithmError {
	return &UnsupportedAlgorithmError{
		Message:   message,
		Algorithm: algorithm,
		Name: "UnsupportedAlgorithmError",
	}
}

func (e *UnsupportedAlgorithmError) Error() string {
	return e.Message
}

func NewRequestError(message string, metadata map[string]interface{}) *RequestError {
	e := NewHTTPError(http.StatusBadRequest, message, metadata)
	return &RequestError{e}
}

func NewAuthorizationError(message string, metadata map[string]interface{}) *AuthorizationError {
	e := NewHTTPError(http.StatusUnauthorized, message, metadata)
	return &AuthorizationError{e}
}

func NewHTTPError(code int, message string, metadata map[string]interface{}) *HTTPError {
	return &HTTPError{
		Code:     code,
		Message:  message,
		Metadata: metadata,
		Name:     http.StatusText(code),
	}
}

func (e *HTTPError) Error() string {
	return string(e.Json())
}

func (e *HTTPError) String() string {
	return fmt.Sprintf("%s (%d): %s, %v", e.Name, e.Code, e.Message, e.Metadata)
}

func (e *HTTPError) Json() []byte {
	val, _ := json.Marshal(e)

	return val
}

func (e *HTTPError) ErrorCode() int {
	return e.Code
}

func (e *HTTPError) ErrorName() string {
	return e.Name
}

func ErrorWriter(err HTTPWrappedError, rw http.ResponseWriter) {
	json := err.Json()

	contentLength := strconv.Itoa(len(json))
	log.WithFields(log.Fields{
		"ContentLength": contentLength,
		"ErrorCode":     err.ErrorCode(),
		"Error":         err,
	}).Infof("Logging %s", err.ErrorName())

	rw.Header().Set("Content-Type", "application/json")
	rw.Header().Set("Content-Length", contentLength)
	rw.WriteHeader(err.ErrorCode())
	rw.Write(json)
}
