package models

type BaseError struct {
	ErrorInterface
	err error
}

func (b *BaseError) SendErrorMessage(err... error) {

}
