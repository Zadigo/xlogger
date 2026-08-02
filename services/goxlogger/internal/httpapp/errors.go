package httpapp

import (
	"net/http"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/utils"
)

type AppErrors struct {
	utils.DefaultErrorResponse
}

func (a *AppErrors) SendErrorMessage(w http.ResponseWriter, err ...error) {

}

func NewErrorHandler() models.ErrorInterface {
	return &AppErrors{}
}
