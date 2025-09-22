package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/soa-team-11/auth-service/api/external"
	"github.com/soa-team-11/auth-service/models"
	"github.com/soa-team-11/auth-service/services"
)

type AuthHandler struct {
	authService     *services.AuthService
	eventService    *external.EventService
	accountsService *services.AccountService
}

func NewAuthHandler() *AuthHandler {
	authService := services.NewAuthService()
	accountsService := services.NewAccountService()
	eventService := external.NewEventService()

	handler := &AuthHandler{
		authService:     authService,
		eventService:    eventService,
		accountsService: accountsService,
	}

	// subscribe na event za kompenzaciju
	handler.eventService.SubscribeCartCreationFailures(func(userID string) error {
		return authService.DeleteUser(userID)
	})

	return handler
}

func (ah *AuthHandler) Routes() chi.Router {
	r := chi.NewRouter()
	r.Use(middleware.AllowContentType("application/json"))

	r.Post("/register", ah.HandleRegister)
	r.Post("/login", ah.HandleLogin)
	r.Get("/list", ah.HandleList)
	r.Patch("/block/{userID}", ah.HandleToggleBlock)

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

	// publish za sagu da se napravi i shopping cart
	ah.eventService.PublishUserRegistered(createdUser.UserID.String())

	w.WriteHeader(http.StatusCreated)
	w.Write(createdUser.ToJSON())
}

func (ah *AuthHandler) HandleList(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	accounts, err := ah.accountsService.ListAccounts()

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

func (ah *AuthHandler) HandleToggleBlock(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userID := chi.URLParam(r, "userID")
	if userID == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"message":"userID is required"}`))
		return
	}
	newStatus, err := ah.accountsService.ToggleBlockUser(userID)
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
