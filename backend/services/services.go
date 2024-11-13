package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type HTTPError interface {
	Data() interface{}
	Msg() string
	StatusCode() int
}

func writeHTTPError(w http.ResponseWriter, e HTTPError) {
	w.WriteHeader(e.StatusCode())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(
		struct {
			Status int         `json:"status"`
			Msg    string      `json:"msg"`
			Data   interface{} `json:"data"`
		}{Status: e.StatusCode(), Msg: e.Msg(), Data: e.Data()})
}

type ServicesHTTPErrorHandler func(func(http.ResponseWriter, *http.Request) error) http.HandlerFunc

func HandleHTTPError(fn func(w http.ResponseWriter, r *http.Request) error) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := fn(w, r)
		if err != nil {
			var se *ServiceError
			switch {
			case errors.As(err, &se):
				writeHTTPError(w, se)
				return
			default:
				writeHTTPError(w, NewInternalServiceError(err))
				return
			}
		}
	}
}

type ServiceError struct {
	data interface{}
	msg  string
	code int
	err  error
}

type ValidationErrMap map[string]ValidationErrData

type ValidationErrData struct {
	Tag   string
	Value interface{}
}

func NewInternalServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusInternalServerError, nil)
}

func NewUnauthenticatedServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusUnauthorized, nil)
}

func NewUnauthorizedServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusForbidden, nil)
}

func NewDataConflictServiceError(err error, msg string) *ServiceError {
	return NewServiceError(err, http.StatusConflict, msg)
}

func NewNotFoundServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusNotFound, nil)
}

func NewBadRequestServiceError(err error) *ServiceError {
	return NewServiceError(err, http.StatusBadRequest, nil)
}

func NewValidationServiceError(err error, data ValidationErrMap) *ServiceError {
	return NewServiceError(err, http.StatusBadRequest, data)
}

func NewServiceError(err error, code int, data interface{}) *ServiceError {
	return &ServiceError{err: err, code: code, msg: http.StatusText(code), data: data}
}

func (se *ServiceError) Data() interface{} {
	return se.data
}

func (se *ServiceError) Msg() string {
	return se.msg
}

func (se *ServiceError) StatusCode() int {
	return se.code
}

func (se *ServiceError) Error() string {
	return fmt.Sprintf("%v: %v", se.code, se.msg)
}

func (se *ServiceError) Unwrap() error {
	return se.err
}
