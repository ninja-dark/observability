package handler

import (
	"net/http"

	"github.com/getsentry/sentry-go"
)

type PanicHandler struct{}




func NewPanicHandler() PanicHandler {
	return PanicHandler{}
}

func (i PanicHandler) Handle(w http.ResponseWriter, r *http.Request) {
	panic("Panic")
}
func (i PanicHandler) Log(w http.ResponseWriter, r *http.Request) {
	//Отправить лог в sentry
	sentry.CaptureMessage("It works!")
}
