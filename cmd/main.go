package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"

	"github.com/soa-team-11/auth-service/api/routers"
	"github.com/soa-team-11/auth-service/utils"
	"github.com/soa-team-11/auth-service/utils/logger"
	"github.com/soa-team-11/auth-service/utils/tracing"
)

func main() {
	logger.Init()

	cleanup := tracing.InitTracer()
	defer cleanup()

	otel.Tracer("auth-service")

	router := routers.Router()
	port := utils.Getenv("PORT", "3001")

	log.Infof("Running services on PORT %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
