package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soa-team-11/auth-service/api/handlers"

	"github.com/soa-team-11/auth-service/middleware"
)

func Router() http.Handler {
	r := chi.NewRouter()
	r.Use(middleware.LogrusMiddleware)

	r.Mount("/auth", handlers.Routes())

	return r
}
