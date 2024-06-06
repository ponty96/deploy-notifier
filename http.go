package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/relistan/apistuff"
	"go.uber.org/zap"
)

type HTTPServer struct {
	LoggingOut io.Writer
}

func (d *HTTPServer) serveHTTP() {
	r := mux.NewRouter()
	r.Use(d.logWrapper)

	// r.HandleFunc("/", HomeHandler)
	// r.HandleFunc("/products", ProductsHandler)
	r.HandleFunc("/health-check", apistuff.HandleHealthOK)

	zap.L().Info("Starting HTTP server on port 8080")
	http.ListenAndServe(":8080", r)
}

func (d *HTTPServer) logWrapper(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip logging for health-check. Log everything else.
		if strings.Contains(r.URL.Path, "health-check") {
			next.ServeHTTP(w, r)
		} else {
			handlers.LoggingHandler(d.LoggingOut, next).ServeHTTP(w, r)
		}
	})
}
