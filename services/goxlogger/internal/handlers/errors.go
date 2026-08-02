package handlers

import (
	"net/http"

	"github.com/Zadigo/goxlogger/internal/models"
	"github.com/Zadigo/goxlogger/internal/utils"
)

type HttpErrors struct {
	utils.DefaultErrorResponse
}

func (a *HttpErrors) BasicErrorMessage(w http.ResponseWriter) {
	a.Detail = "An error occurred while processing your request."
	a.Message = "An unexpected error occurred."
	a.SendErrorMessage(w)
}

func (a *HttpErrors) InvalidFileId(w http.ResponseWriter) {
	a.Detail = "The provided file ID is invalid."
	a.Message = "Invalid file ID."
	a.SendErrorMessage(w)
}

func (a *HttpErrors) FailedToGetLogs(w http.ResponseWriter, err ...error) {
	a.Detail = "Failed to retrieve logs for the provided file ID."
	a.Message = "Failed to get logs."
	a.SendErrorMessage(w, err...)
}

func NewErrorHandler() models.ErrorInterface {
	return &HttpErrors{}
}
