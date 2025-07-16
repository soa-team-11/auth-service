package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/soa-team-11/auth-service/models"
	"github.com/soa-team-11/auth-service/services"
)

var (
	authService = services.NewAuthService()
)

func Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/register", HandleRegister)
	r.Post("/login", nil)

	return r
}

func HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	b, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	createdUser, err := authService.Register(user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"message":"%s"}`, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(createdUser.ToJSON())
}
