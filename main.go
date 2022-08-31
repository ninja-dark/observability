package main

import (
	"elastic/handler"
	"elastic/l"
	sentrylog "elastic/sentry"
	"elastic/store"
	"log"

	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
)

// Переписать не на Martini
func main() {

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	defer func() { _ = logger.Sync() }()

	tracer, closer := l.InitJaeger("goweb", logger)
	defer closer.Close()

	sentrylog.SentryLog()

	//Initialize Stores
	articleStore, err := store.NewArticleStore(logger, tracer)
	parseErr(err)
	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore, logger, tracer)
	panicHandler := handler.PanicHandler{}
	//Initialize Router
	r := chi.NewRouter()
	//Routes
	r.Get("/article/id/:id", articleHandler.Id)
	r.Post("/article/add", articleHandler.Add)
	r.Post("/article/search", articleHandler.Search)
	 r.Get("/", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello World!"))
    })

	r.Get("/panic", panicHandler.Handle)
	r.Post("/log/add", panicHandler.Log)
	http.ListenAndServe(":8080", r)
}

func parseErr(err error) {
	if err != nil {
		l.F(err)
	}
	l.L("Application started")
}
