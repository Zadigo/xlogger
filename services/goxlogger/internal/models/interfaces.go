package models

import "net/http"

type ErrorInterface interface {
	LogErrorMessage(err ...error)
	SendErrorMessage(w http.ResponseWriter, err ...error)
}
