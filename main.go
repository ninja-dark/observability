package main

import (
	"html/template"
	"net/http"
	"observability/middleware"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	r := mux.NewRouter()

	metricsMiddleware := middleware.NewMetricsMiddleware()

	r.Handle("/metrics", promhttp.Handler())
	r.HandleFunc("/alert", alertHandler).Methods(http.MethodGet)
	r.HandleFunc("/simple", simpleHandler).Methods(http.MethodPost)
	r.HandleFunc("/hard", hardHandler).Methods(http.MethodPut)

	r.Use(metricsMiddleware.Metrics)

	http.ListenAndServe(":8080", r)
}

func alertHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Alert"))
}

func hardHandler(w http.ResponseWriter, r *http.Request) {
	path := filepath.Join("static", "hard.html")
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte("Simple"))
}
