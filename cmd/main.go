package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/soa-team-11/auth-service/api/routers"
	"github.com/soa-team-11/auth-service/utils"
)

func main() {
	router := routers.Router()

	port := utils.Getenv("PORT", "3001")

	log.Infof("Running services on PORT %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
