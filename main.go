package main

import (
	"elastic/handler"
	"elastic/l"
	sentrylog "elastic/sentry"
	"elastic/store"

	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

// Переписать не на Martini
func main() {

	sentrylog.SentryLog()

	//Initialize Stores
	articleStore, err := store.NewArticleStore()
	parseErr(err)
	//Initialize Handlers
	articleHandler := handler.NewArticleHandler(articleStore)
	panicHandler := handler.PanicHandler{}
	//Initialize Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
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
	l.Log.Log("Application started")

	
}
