package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/soa-team-11/auth-service/api/external"
	"github.com/soa-team-11/auth-service/models"
	"github.com/soa-team-11/auth-service/services"
)

type AuthHandler struct {
	authService  *services.AuthService
	eventService *external.EventService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{authService: services.NewAuthService(), eventService: external.NewEventService()}
}

func (ah *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/register", ah.HandleRegister)
	r.Post("/login", ah.HandleLogin)

	return r
}

func (ah *AuthHandler) HandleLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	request_data := struct {
		Username string `json:"username" bson:"username"`
		Password string `json:"password,omitempty" bson:"password"`
	}{}

	err = json.Unmarshal(b, &request_data)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	login, err := ah.authService.Login(request_data.Username, request_data.Password)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(login)
}

func (ah *AuthHandler) HandleRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	defer r.Body.Close()

	b, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	var user models.User
	err = json.Unmarshal(b, &user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	createdUser, err := ah.authService.Register(user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"message": err.Error()})
		return
	}

	ah.eventService.PublishUserRegistered(createdUser.UserID.String())

	w.WriteHeader(http.StatusCreated)
	w.Write(createdUser.ToJSON())
}
