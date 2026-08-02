package models

type ErrorInterface interface {
	SendErrorMessage(err... error)
}
