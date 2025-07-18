package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/soa-team-11/auth-service/services"
)

type AccountsHandler struct {
	accountService *services.AccountService
}

func NewAccountsHandler() *AccountsHandler {
	return &AccountsHandler{accountService: services.NewAccountService()}
}

func (ah *AccountsHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))

	r.Get("/list", ah.HandleListAccounts)

	return r
}

func (ah *AccountsHandler) HandleListAccounts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accounts, err := ah.accountService.ListAccounts()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	accountsJSON, err := json.Marshal(accounts)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(accountsJSON)
}
