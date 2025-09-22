package routers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/soa-team-11/auth-service/api/handlers"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"

	"github.com/soa-team-11/auth-service/middleware"
)

var (
	authHandler = handlers.NewAuthHandler()
)

func Router() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.LogrusMiddleware)
	r.Use(otelhttp.NewMiddleware("auth-service"))

	r.Mount("/auth", authHandler.Routes())

	return r
}
