package handler

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/opentracing/opentracing-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type PanicHandler struct{
	logger *zap.Logger
	tracer opentracing.Tracer
}


func NewPanicHandler(logger *zap.Logger, tracer opentracing.Tracer) PanicHandler {
	return PanicHandler{
		logger: logger,
		tracer: tracer,
	}
}

func (i PanicHandler) Handle(w http.ResponseWriter, r *http.Request) {
	
	i.logger.Info("PanicHandler handle called", zap.Field{Key: "method", String: r.Method, Type:zapcore.StringType})
	panic("Panic")
}
func (i PanicHandler) Log(w http.ResponseWriter, r *http.Request) {
	//Отправить лог в sentry
	sentry.CaptureMessage("It works!")

	writeJsonResponse(w, http.StatusOK, "sent log to sentry")
}
