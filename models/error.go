package models

import "net/http"

type Errors struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Errors  interface{} `json:"errors"`
}

func (ie *Errors) Error() string {
	return ie.Message
}

func NewErrors(code int, msg string, data interface{}) *Errors {
	return &Errors{
		Code:    code,
		Message: msg,
		Errors:  data,
	}
}

func NewBadRequest(data interface{}) *Errors {
	return &Errors{
		Code:    http.StatusBadRequest,
		Message: "bad request",
		Errors:  data,
	}
}
func NewBadRequestID() *Errors {
	return &Errors{
		Code:    http.StatusBadRequest,
		Message: "bad request",
		Errors:  "Error parameter id",
	}
}

func NewFailedToConvertData() *Errors {
	return &Errors{
		Code:    http.StatusBadRequest,
		Message: "Internal server error",
		Errors:  "CANNOT CONVERT ID PARAMS",
	}
}

func NewInternalServerError(data interface{}) *Errors {
	return &Errors{
		Code:    http.StatusInternalServerError,
		Message: "INTERNAL SERVER ERROR",
		Errors:  data,
	}
}

func NewNotFound(data interface{}) *Errors {
	return &Errors{
		Code:    http.StatusInternalServerError,
		Message: "not found",
		Errors:  data,
	}
}
