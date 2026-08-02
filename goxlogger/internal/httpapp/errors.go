package httpapp

import "github.com/Zadigo/goxlogger/internal/models"

type AppErrors struct{}

func (a *AppErrors) SendErrorMessage(err ...error) {

}

func NewErrorHandler() models.ErrorInterface {
	return &AppErrors{}
}
