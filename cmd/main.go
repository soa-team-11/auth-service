package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/soa-team-11/auth-service/config"
)

func main() {
	cfg := config.LoadConfig()

	r := chi.NewRouter()
	r.Use(middleware.Logger)

	log.Printf("Starting server on port %s", cfg.Port)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
