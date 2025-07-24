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
	r.Patch("/block/{userID}", ah.HandleToggleBlockUser)

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

// HandleToggleBlockUser toggles the blocked status of a user and returns the new status as JSON
// Endpoint: PATCH /block/{userID}
func (ah *AccountsHandler) HandleToggleBlockUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"userID is required"}`))
		return
	}
	newStatus, err := ah.accountService.ToggleBlockUser(userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf(`{"message":"%s"}`, err.Error())))
		return
	}
	response := struct {
		UserID  string `json:"user_id"`
		Blocked bool   `json:"blocked"`
	}{
		UserID:  userID,
		Blocked: newStatus,
	}
	jsonResp, _ := json.Marshal(response)
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResp)
}
