package handlers

import "github.com/Zadigo/goxlogger/internal/models"

type HttpErrors struct{}

func (a *HttpErrors) SendErrorMessage(err ...error) {

}

func NewErrorHandler() models.ErrorInterface {
	return &HttpErrors{}
}
